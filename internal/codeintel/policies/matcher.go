package policies

import (
	"context"
	"fmt"
	"time"

	"github.com/gobwas/glob"

	"github.com/sourcegraph/sourcegraph/internal/api"
	"github.com/sourcegraph/sourcegraph/internal/codeintel/policies/shared"
	"github.com/sourcegraph/sourcegraph/internal/gitserver"
	"github.com/sourcegraph/sourcegraph/internal/gitserver/gitdomain"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

type Matcher struct {
	gitserverClient           gitserver.Client
	extractor                 Extractor
	includeTipOfDefaultBranch bool
	filterByCreatedDate       bool
}

// PolicyMatch indicates the name of the matching branch or tag associated with some commit. The policy
// identifier field is set unless the policy match exists due to a `includeTipOfDefaultBranch` match. The
// policy duration field is set if the matching policy specifies a duration.
type PolicyMatch struct {
	Name           string
	PolicyID       *int
	PolicyDuration *time.Duration
	CommittedAt    time.Time
}

func NewMatcher(
	gitserverClient gitserver.Client,
	extractor Extractor,
	includeTipOfDefaultBranch bool,
	filterByCreatedDate bool,
) *Matcher {
	return &Matcher{
		gitserverClient:           gitserverClient,
		extractor:                 extractor,
		includeTipOfDefaultBranch: includeTipOfDefaultBranch,
		filterByCreatedDate:       filterByCreatedDate,
	}
}

// CommitsDescribedByPolicy returns a map from commits within the given repository to a set of policy matches
// with respect to the given policies.
//
// If includeTipOfDefaultBranch is true, then there will exist a match for the tip default branch with a nil
// policy identifier, even if no policies are supplied. This is set to true for data retention but not for
// auto-indexing.
//
// If filterByCreatedDate is true, then commits that are older than the matching policy duration will be
// filtered out. If false, policy duration is not considered. This is set to true for auto-indexing, but false
// for data retention as we need to compare the policy duration against the associated upload date, not the
// commit date.
//
// A subset of all commits can be returned by passing in any number of commit revhash strings.
func (m *Matcher) CommitsDescribedByPolicy(ctx context.Context, repositoryID int, repoName api.RepoName, policies []shared.ConfigurationPolicy, now time.Time, filterCommits ...string) (map[string][]PolicyMatch, error) {
	if len(policies) == 0 && !m.includeTipOfDefaultBranch {
		return nil, nil
	}

	patterns, err := compilePatterns(policies)
	if err != nil {
		return nil, err
	}

	// mutable context
	mContext := matcherContext{
		repositoryID:   repositoryID,
		repo:           repoName,
		policies:       policies,
		patterns:       patterns,
		commitMap:      map[string][]PolicyMatch{},
		branchRequests: map[string]branchRequestMeta{},
	}

	opt := gitserver.ListRefsOpts{
		HeadsOnly: true,
		TagsOnly:  true,
	}
	for _, c := range filterCommits {
		opt.PointsAtCommit = append(opt.PointsAtCommit, api.CommitID(c))
	}

	refs, err := m.gitserverClient.ListRefs(ctx, repoName, opt)
	if err != nil {
		return nil, errors.Wrap(err, "gitserver.ListRefs")
	}

	for _, ref := range refs {
		switch ref.Type {
		case gitdomain.RefTypeTag:
			// Match tagged commits
			m.matchTaggedCommits(mContext, string(ref.CommitID), ref, now)

		case gitdomain.RefTypeBranch:
			// Match tips of branches
			m.matchBranchHeads(mContext, string(ref.CommitID), ref, now)
		}
	}

	// Match commits on branches but not at tip
	if err := m.matchCommitsOnBranch(ctx, mContext, now); err != nil {
		return nil, err
	}

	// Match comments via rev-parse
	if err := m.matchCommitPolicies(ctx, mContext, now); err != nil {
		return nil, err
	}

	return mContext.commitMap, nil
}

type matcherContext struct {
	// repositoryID is the repository identifier used in requests to gitserver.
	repositoryID int

	repo api.RepoName

	// policies is the full set (global and repository-specific) policies that apply to the given repository.
	policies []shared.ConfigurationPolicy

	// patterns holds a compiled glob of the pattern from each non-commit type policy.
	patterns map[string]glob.Glob

	// commitMap stores matching policies for each commit in the given repository.
	commitMap map[string][]PolicyMatch

	// branchRequests holds metadata about the additional requests we need to make to gitserver to determine
	// the set of commits that are an ancestor of a branch head (but not an ancestor of the default branch).
	// These commits are "contained" within in the intermediate commits composing a logical branch in the git
	// graph. As multiple policies can match the same branch, we store it in a map to ensure that we make only
	// one request per branch.
	branchRequests map[string]branchRequestMeta
}

type branchRequestMeta struct {
	isDefaultBranch     bool
	commitID            string // commit hash of the tip of the branch
	policyDurationByIDs map[int]*time.Duration
}

// matchTaggedCommits determines if the given commit (described by the tag-type ref given description) matches any tag-type
// policies. For each match, a commit/policy pair will be added to the given context.
func (m *Matcher) matchTaggedCommits(context matcherContext, commit string, ref gitdomain.Ref, now time.Time) {
	visitor := func(policy shared.ConfigurationPolicy) {
		policyDuration, _ := m.extractor(policy)

		context.commitMap[commit] = append(context.commitMap[commit], PolicyMatch{
			Name:           ref.ShortName,
			PolicyID:       &policy.ID,
			PolicyDuration: policyDuration,
			CommittedAt:    ref.CreatedDate,
		})
	}

	m.forEachMatchingPolicy(context, ref, shared.GitObjectTypeTag, visitor, now)
}

// matchBranchHeads determines if the given commit (described by the branch-type ref given description) matches any branch-type
// policies. For each match, a commit/policy pair will be added to the given context. This method also adds matches for the tip
// of the default branch (if configured to do so), and adds bookkeeping metadata to the context's branchRequests field when a
// matching policy's intermediate commits should be checked.
func (m *Matcher) matchBranchHeads(context matcherContext, commit string, ref gitdomain.Ref, now time.Time) {
	if ref.IsHead && m.includeTipOfDefaultBranch {
		// Add a match with no associated policy for the tip of the default branch
		context.commitMap[commit] = append(context.commitMap[commit], PolicyMatch{
			Name:           ref.ShortName,
			PolicyID:       nil,
			PolicyDuration: nil,
			CommittedAt:    ref.CreatedDate,
		})
	}

	visitor := func(policy shared.ConfigurationPolicy) {
		policyDuration, _ := m.extractor(policy)

		context.commitMap[commit] = append(context.commitMap[commit], PolicyMatch{
			Name:           ref.ShortName,
			PolicyID:       &policy.ID,
			PolicyDuration: policyDuration,
			CommittedAt:    ref.CreatedDate,
		})

		// Build requests to be made in batch via the matchCommitsOnBranch method
		if policyDuration, includeIntermediateCommits := m.extractor(policy); includeIntermediateCommits {
			meta, ok := context.branchRequests[ref.ShortName]
			if !ok {
				meta.policyDurationByIDs = map[int]*time.Duration{}
			}

			meta.policyDurationByIDs[policy.ID] = policyDuration
			meta.isDefaultBranch = meta.isDefaultBranch || ref.IsHead
			meta.commitID = commit
			context.branchRequests[ref.ShortName] = meta
		}
	}

	m.forEachMatchingPolicy(context, ref, shared.GitObjectTypeTree, visitor, now)
}

// matchCommitsOnBranch makes a request for commits belonging to any branch matching a branch-type
// policy that also includes intermediate commits. This method uses the requests queued by the
// matchBranchHeads method. A commit/policy pair will be added to the given context for each commit
// of appropriate age existing on a matched branch.
func (m *Matcher) matchCommitsOnBranch(ctx context.Context, context matcherContext, now time.Time) error {
	for branchName, branchRequestMeta := range context.branchRequests {
		maxCommitAge := getMaxAge(branchRequestMeta.policyDurationByIDs, now)

		if !m.filterByCreatedDate {
			// Do not filter out any commits by date
			maxCommitAge = nil
		}

		commitDates, err := commitsUniqueToBranch(
			ctx,
			m.gitserverClient,
			context.repo,
			api.CommitID(branchRequestMeta.commitID),
			branchRequestMeta.isDefaultBranch,
			maxCommitAge,
		)
		if err != nil {
			return errors.Wrap(err, "gitserver.CommitsUniqueToBranch")
		}

		for commit, commitDate := range commitDates {
		policyLoop:
			for policyID, policyDuration := range branchRequestMeta.policyDurationByIDs {
				for _, match := range context.commitMap[string(commit)] {
					if match.PolicyID != nil && *match.PolicyID == policyID {
						// Skip duplicates (can happen at head of branch)
						continue policyLoop
					}
				}

				if m.filterByCreatedDate && policyDuration != nil && now.Sub(commitDate) > *policyDuration {
					// Policy duration was less than max age and re-check failed
					continue policyLoop
				}

				context.commitMap[string(commit)] = append(context.commitMap[string(commit)], PolicyMatch{
					Name:           branchName,
					PolicyID:       &policyID,
					PolicyDuration: policyDuration,
					CommittedAt:    commitDate,
				})
			}
		}
	}

	return nil
}

// commitsUniqueToBranch returns a map from commits that exist on a particular
// branch in the given repository to their committer date. This set of commits is
// determined by listing `{branchName} ^HEAD`, which is interpreted as: all
// commits on {branchName} not also on the tip of the default branch. If the
// supplied branch name is the default branch, then this method instead returns
// all commits reachable from HEAD.
func commitsUniqueToBranch(ctx context.Context, gitserverClient gitserver.Client, repo api.RepoName, commitID api.CommitID, isDefaultBranch bool, maxCommitAge *time.Time) (map[api.CommitID]time.Time, error) {
	rng := "HEAD"
	if !isDefaultBranch {
		rng = fmt.Sprintf("HEAD..%s", commitID)
	}

	var after time.Time
	if maxCommitAge != nil {
		after = *maxCommitAge
	}

	commits, err := gitserverClient.Commits(ctx, repo, gitserver.CommitsOptions{
		Ranges: []string{rng},
		After:  after,
	})
	if err != nil {
		return nil, err
	}

	commitMap := make(map[api.CommitID]time.Time)
	for _, commit := range commits {
		commitMap[commit.ID] = commit.Committer.Date
	}

	return commitMap, nil

}

// matchCommitPolicies compares the each commit-type policy pattern as a rev-like against the target
// repository in gitserver. For each match, a commit/policy pair will be added to the given context.
func (m *Matcher) matchCommitPolicies(ctx context.Context, context matcherContext, now time.Time) error {
	for _, policy := range context.policies {
		if policy.Type == shared.GitObjectTypeCommit {
			commit, err := m.gitserverClient.GetCommit(ctx, context.repo, api.CommitID(policy.Pattern))
			if err != nil {
				if errors.HasType(err, &gitdomain.RevisionNotFoundError{}) {
					continue
				}
				return err
			}

			policyDuration, _ := m.extractor(policy)

			if m.filterByCreatedDate && policyDuration != nil && now.Sub(commit.Committer.Date) > *policyDuration {
				continue
			}

			context.commitMap[policy.Pattern] = append(context.commitMap[policy.Pattern], PolicyMatch{
				Name:           string(commit.ID),
				PolicyID:       &policy.ID,
				PolicyDuration: policyDuration,
				CommittedAt:    commit.Committer.Date,
			})
		}
	}

	return nil
}

func (m *Matcher) forEachMatchingPolicy(context matcherContext, ref gitdomain.Ref, targetObjectType shared.GitObjectType, f func(policy shared.ConfigurationPolicy), now time.Time) {
	for _, policy := range context.policies {
		if policy.Type == targetObjectType && m.policyMatchesRef(context, policy, ref, now) {
			f(policy)
		}
	}
}

func (m *Matcher) policyMatchesRef(context matcherContext, policy shared.ConfigurationPolicy, ref gitdomain.Ref, now time.Time) bool {
	if !context.patterns[policy.Pattern].Match(ref.ShortName) {
		// Name doesn't match policy's pattern
		return false
	}

	if policyDuration, _ := m.extractor(policy); m.filterByCreatedDate && policyDuration != nil && now.Sub(ref.CreatedDate) > *policyDuration {
		// Policy is not unbounded, we are filtering by commit date, commit is too old
		return false
	}

	return true
}

// compilePatterns constructs a map from patterns in each given policy to a compiled glob object used
// to match to commits, branch names, and tag names. If there are multiple policies with the same pattern,
// the pattern is compiled only once.
func compilePatterns(policies []shared.ConfigurationPolicy) (map[string]glob.Glob, error) {
	patterns := make(map[string]glob.Glob, len(policies))
	for _, policy := range policies {
		if _, ok := patterns[policy.Pattern]; ok || policy.Type == shared.GitObjectTypeCommit {
			continue
		}

		pattern, err := glob.Compile(policy.Pattern)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to compile glob pattern `%s` in configuration policy %d", policy.Pattern, policy.ID))
		}

		patterns[policy.Pattern] = pattern
	}

	return patterns, nil
}

func getMaxAge(policyDurationByIDs map[int]*time.Duration, now time.Time) *time.Time {
	var maxDuration *time.Duration
	for _, duration := range policyDurationByIDs {
		if duration == nil {
			// If any duration is nil, the policy is unbounded
			return nil
		}
		if maxDuration == nil || *maxDuration < *duration {
			maxDuration = duration
		}
	}
	if maxDuration == nil {
		return nil
	}

	maxAge := now.Add(-*maxDuration)
	return &maxAge
}
