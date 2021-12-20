import React, { Suspense, useCallback, useEffect, useMemo } from 'react'
import { Redirect, Route, RouteComponentProps, Switch, matchPath } from 'react-router'
import { Observable } from 'rxjs'

import { ResizablePanel } from '@sourcegraph/branded/src/components/panel/Panel'
import { LoadingSpinner } from '@sourcegraph/react-loading-spinner'
import { ActivationProps } from '@sourcegraph/shared/src/components/activation/Activation'
import { FetchFileParameters } from '@sourcegraph/shared/src/components/CodeExcerpt'
import { ExtensionsControllerProps } from '@sourcegraph/shared/src/extensions/controller'
import * as GQL from '@sourcegraph/shared/src/graphql/schema'
import { PlatformContextProps } from '@sourcegraph/shared/src/platform/context'
import { SettingsCascadeProps } from '@sourcegraph/shared/src/settings/settings'
import { TelemetryProps } from '@sourcegraph/shared/src/telemetry/telemetryService'
import { parseQueryAndHash } from '@sourcegraph/shared/src/util/url'
import { useObservable } from '@sourcegraph/shared/src/util/useObservable'

import { AuthenticatedUser, authRequired as authRequiredObservable } from './auth'
import { BatchChangesProps } from './batches'
import { CodeIntelligenceProps } from './codeintel'
import { communitySearchContextsRoutes } from './communitySearchContexts/routes'
import { AppRouterContainer } from './components/AppRouterContainer'
import { useBreadcrumbs } from './components/Breadcrumbs'
import { ErrorBoundary } from './components/ErrorBoundary'
import { useScrollToLocationHash } from './components/useScrollToLocationHash'
import { GlobalContributions } from './contributions'
import { ExtensionAreaRoute } from './extensions/extension/ExtensionArea'
import { ExtensionAreaHeaderNavItem } from './extensions/extension/ExtensionAreaHeader'
import { ExtensionsAreaRoute } from './extensions/ExtensionsArea'
import { ExtensionsAreaHeaderActionButton } from './extensions/ExtensionsAreaHeader'
import { FeatureFlagProps } from './featureFlags/featureFlags'
import { GlobalAlerts } from './global/GlobalAlerts'
import { GlobalDebug } from './global/GlobalDebug'
import { CodeInsightsContextProps, CodeInsightsProps } from './insights/types'
import { KeyboardShortcutsProps, KEYBOARD_SHORTCUT_SHOW_HELP } from './keyboardShortcuts/keyboardShortcuts'
import { KeyboardShortcutsHelp } from './keyboardShortcuts/KeyboardShortcutsHelp'
import styles from './Layout.module.scss'
import { SurveyToast } from './marketing/SurveyToast'
import { GlobalNavbar } from './nav/GlobalNavbar'
import { useExtensionAlertAnimation } from './nav/UserNavItem'
import { OrgAreaRoute } from './org/area/OrgArea'
import { OrgAreaHeaderNavItem } from './org/area/OrgHeader'
import { fetchHighlightedFileLineRanges } from './repo/backend'
import { RepoContainerRoute } from './repo/RepoContainer'
import { RepoHeaderActionButton } from './repo/RepoHeader'
import { RepoRevisionContainerRoute } from './repo/RepoRevisionContainer'
import { RepoSettingsAreaRoute } from './repo/settings/RepoSettingsArea'
import { RepoSettingsSideBarGroup } from './repo/settings/RepoSettingsSidebar'
import { LayoutRouteProps, LayoutRouteComponentProps } from './routes'
import { PageRoutes, EnterprisePageRoutes } from './routes.constants'
import { Settings } from './schema/settings.schema'
import {
    parseSearchURLQuery,
    PatternTypeProps,
    HomePanelsProps,
    SearchStreamingProps,
    parseSearchURL,
    SearchContextProps,
    getGlobalSearchContextFilter,
} from './search'
import { SiteAdminAreaRoute } from './site-admin/SiteAdminArea'
import { SiteAdminSideBarGroups } from './site-admin/SiteAdminSidebar'
import { setQueryStateFromURL } from './stores'
import { useThemeProps } from './theme'
import { UserAreaRoute } from './user/area/UserArea'
import { UserAreaHeaderNavItem } from './user/area/UserAreaHeader'
import { UserSettingsAreaRoute } from './user/settings/UserSettingsArea'
import { UserSettingsSidebarItems } from './user/settings/UserSettingsSidebar'
import { isMacPlatform, UserExternalServicesOrRepositoriesUpdateProps } from './util'
import { parseBrowserRepoURL } from './util/url'

export interface LayoutProps
    extends RouteComponentProps<{}>,
        SettingsCascadeProps<Settings>,
        PlatformContextProps,
        ExtensionsControllerProps,
        KeyboardShortcutsProps,
        TelemetryProps,
        ActivationProps,
        PatternTypeProps,
        SearchContextProps,
        HomePanelsProps,
        SearchStreamingProps,
        UserExternalServicesOrRepositoriesUpdateProps,
        CodeIntelligenceProps,
        BatchChangesProps,
        CodeInsightsProps,
        CodeInsightsContextProps,
        FeatureFlagProps {
    extensionAreaRoutes: readonly ExtensionAreaRoute[]
    extensionAreaHeaderNavItems: readonly ExtensionAreaHeaderNavItem[]
    extensionsAreaRoutes: readonly ExtensionsAreaRoute[]
    extensionsAreaHeaderActionButtons: readonly ExtensionsAreaHeaderActionButton[]
    siteAdminAreaRoutes: readonly SiteAdminAreaRoute[]
    siteAdminSideBarGroups: SiteAdminSideBarGroups
    siteAdminOverviewComponents: readonly React.ComponentType[]
    userAreaHeaderNavItems: readonly UserAreaHeaderNavItem[]
    userAreaRoutes: readonly UserAreaRoute[]
    userSettingsSideBarItems: UserSettingsSidebarItems
    userSettingsAreaRoutes: readonly UserSettingsAreaRoute[]
    orgAreaHeaderNavItems: readonly OrgAreaHeaderNavItem[]
    orgAreaRoutes: readonly OrgAreaRoute[]
    repoContainerRoutes: readonly RepoContainerRoute[]
    repoRevisionContainerRoutes: readonly RepoRevisionContainerRoute[]
    repoHeaderActionButtons: readonly RepoHeaderActionButton[]
    repoSettingsAreaRoutes: readonly RepoSettingsAreaRoute[]
    repoSettingsSidebarGroups: readonly RepoSettingsSideBarGroup[]
    routes: readonly LayoutRouteProps<any>[]

    authenticatedUser: AuthenticatedUser | null

    /**
     * The subject GraphQL node ID of the viewer, which is used to look up the viewer's settings. This is either
     * the site's GraphQL node ID (for anonymous users) or the authenticated user's GraphQL node ID.
     */
    viewerSubject: Pick<GQL.ISettingsSubject, 'id' | 'viewerCanAdminister'>

    // Search
    fetchHighlightedFileLineRanges: (parameters: FetchFileParameters, force?: boolean) => Observable<string[][]>

    globbing: boolean
    isSourcegraphDotCom: boolean
    fetchSavedSearches: () => Observable<GQL.ISavedSearch[]>
    children?: never
}

export const Layout: React.FunctionComponent<LayoutProps> = props => {
    const routeMatch = props.routes.find(({ path, exact }) => matchPath(props.location.pathname, { path, exact }))?.path
    const isSearchRelatedPage = (routeMatch === '/:repoRevAndRest+' || routeMatch?.startsWith('/search')) ?? false
    const minimalNavLinks = routeMatch === '/cncf'
    const isSearchHomepage = props.location.pathname === '/search' && !parseSearchURLQuery(props.location.search)
    const isSearchConsolePage = routeMatch?.startsWith('/search/console')
    const isSearchNotebookPage = routeMatch?.startsWith('/search/notebook')
    const isRepositoryRelatedPage = routeMatch === '/:repoRevAndRest+' ?? false

    // Update patternType, caseSensitivity, and selectedSearchContextSpec based on current URL
    const {
        history,
        patternType: currentPatternType,
        selectedSearchContextSpec,
        location,
        setPatternType,
        setSelectedSearchContextSpec,
    } = props

    useEffect(() => setQueryStateFromURL(location.search), [location.search])

    const { query = '', patternType } = useMemo(() => parseSearchURL(location.search), [location.search])

    const searchContextSpec = useMemo(() => getGlobalSearchContextFilter(query)?.spec, [query])

    useEffect(() => {
        // Only override filters from URL if there is a search query
        if (query) {
            if (patternType && patternType !== currentPatternType) {
                setPatternType(patternType)
            }

            if (searchContextSpec && searchContextSpec !== selectedSearchContextSpec) {
                setSelectedSearchContextSpec(searchContextSpec)
            }
        }
    }, [
        history,
        currentPatternType,
        selectedSearchContextSpec,
        patternType,
        query,
        setPatternType,
        setSelectedSearchContextSpec,
        searchContextSpec,
    ])

    const communitySearchContextPaths = communitySearchContextsRoutes.map(route => route.path)
    const isCommunitySearchContextPage = communitySearchContextPaths.includes(props.location.pathname)

    // TODO add a component layer as the parent of the Layout component rendering "top-level" routes that do not render the navbar,
    // so that Layout can always render the navbar.
    const needsSiteInit = window.context?.needsSiteInit
    const isSiteInit = props.location.pathname === PageRoutes.SiteAdminInit
    const isSignInOrUp =
        props.location.pathname === PageRoutes.SignIn ||
        props.location.pathname === PageRoutes.SignUp ||
        props.location.pathname === PageRoutes.PasswordReset ||
        props.location.pathname === PageRoutes.Welcome

    // TODO Change this behavior when we have global focus management system
    // Need to know this for disable autofocus on nav search input
    // and preserve autofocus for first textarea at survey page, creation UI etc.
    const isSearchAutoFocusRequired = routeMatch === PageRoutes.Survey || routeMatch === EnterprisePageRoutes.Insights

    const authRequired = useObservable(authRequiredObservable)

    const themeProps = useThemeProps()

    const breadcrumbProps = useBreadcrumbs()

    // Control browser extension discoverability animation here.
    // `Layout` is the lowest common ancestor of `UserNavItem` (target) and `RepoContainer` (trigger)
    const { isExtensionAlertAnimating, startExtensionAlertAnimation } = useExtensionAlertAnimation()
    const onExtensionAlertDismissed = useCallback(() => {
        startExtensionAlertAnimation()
    }, [startExtensionAlertAnimation])

    useScrollToLocationHash(props.location)
    // Remove trailing slash (which is never valid in any of our URLs).
    if (props.location.pathname !== '/' && props.location.pathname.endsWith('/')) {
        return <Redirect to={{ ...props.location, pathname: props.location.pathname.slice(0, -1) }} />
    }

    const context: LayoutRouteComponentProps<any> = {
        ...props,
        ...themeProps,
        ...breadcrumbProps,
        onExtensionAlertDismissed,
        isMacPlatform,
        parsedSearchQuery: query,
    }

    return (
        <div className={styles.layout}>
            <KeyboardShortcutsHelp
                keyboardShortcutForShow={KEYBOARD_SHORTCUT_SHOW_HELP}
                keyboardShortcuts={props.keyboardShortcuts}
            />
            <GlobalAlerts authenticatedUser={props.authenticatedUser} settingsCascade={props.settingsCascade} />
            {!isSiteInit && <SurveyToast />}
            {!isSiteInit && !isSignInOrUp && (
                <GlobalNavbar
                    {...props}
                    {...themeProps}
                    parsedSearchQuery={query}
                    authRequired={!!authRequired}
                    showSearchBox={
                        isSearchRelatedPage &&
                        !isSearchHomepage &&
                        !isCommunitySearchContextPage &&
                        !isSearchConsolePage &&
                        !isSearchNotebookPage
                    }
                    variant={
                        isSearchHomepage
                            ? 'low-profile'
                            : isCommunitySearchContextPage
                            ? 'low-profile-with-logo'
                            : 'default'
                    }
                    minimalNavLinks={minimalNavLinks}
                    isSearchAutoFocusRequired={!isSearchAutoFocusRequired}
                    isExtensionAlertAnimating={isExtensionAlertAnimating}
                    isRepositoryRelatedPage={isRepositoryRelatedPage}
                />
            )}
            {needsSiteInit && !isSiteInit && <Redirect to="/site-admin/init" />}
            <ErrorBoundary location={props.location}>
                <Suspense
                    fallback={
                        <div className="flex flex-1">
                            <LoadingSpinner className="icon-inline m-2" />
                        </div>
                    }
                >
                    <Switch>
                        {props.routes.map(
                            ({ render, condition = () => true, ...route }) =>
                                condition(context) && (
                                    <Route
                                        {...route}
                                        key="hardcoded-key" // see https://github.com/ReactTraining/react-router/issues/4578#issuecomment-334489490
                                        component={undefined}
                                        render={routeComponentProps => (
                                            <AppRouterContainer>
                                                {render({ ...context, ...routeComponentProps })}
                                            </AppRouterContainer>
                                        )}
                                    />
                                )
                        )}
                    </Switch>
                </Suspense>
            </ErrorBoundary>
            {parseQueryAndHash(props.location.search, props.location.hash).viewState &&
                props.location.pathname !== PageRoutes.SignIn && (
                    <ResizablePanel
                        {...props}
                        {...themeProps}
                        repoName={`git://${parseBrowserRepoURL(props.location.pathname).repoName}`}
                        fetchHighlightedFileLineRanges={fetchHighlightedFileLineRanges}
                    />
                )}
            <GlobalContributions
                key={3}
                extensionsController={props.extensionsController}
                platformContext={props.platformContext}
                history={props.history}
            />
            <GlobalDebug {...props} />
        </div>
    )
}
