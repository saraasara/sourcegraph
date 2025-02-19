<svelte:options immutable />

<script lang="ts">
    import { mdiFolderArrowUpOutline, mdiFolderOpenOutline, mdiFolderOutline } from '@mdi/js'

    import { goto } from '$app/navigation'
    import Icon from '$lib/Icon.svelte'
    import Popover from '$lib/Popover.svelte'
    import { type FileTreeProvider, NODE_LIMIT, type TreeEntry } from '$lib/repo/api/tree'
    import FileIcon from '$lib/repo/FileIcon.svelte'
    import FilePopover, { fetchPopoverData } from '$lib/repo/filePopover/FilePopover.svelte'
    import { getSidebarFileTreeStateForRepo } from '$lib/repo/stores'
    import { replaceRevisionInURL } from '$lib/shared'
    import TreeView, { setTreeContext } from '$lib/TreeView.svelte'
    import { createForwardStore } from '$lib/utils'
    import { Alert } from '$lib/wildcard'

    export let repoName: string
    export let treeProvider: FileTreeProvider
    export let selectedPath: string
    export let revision: string

    /**
     * Returns the corresponding icon for `entry`
     */
    function getDirectoryIconPath(entry: TreeEntry, open: boolean) {
        if (entry === treeRoot) {
            return mdiFolderArrowUpOutline
        }
        return open ? mdiFolderOpenOutline : mdiFolderOutline
    }

    /**
     * Navigates to the tree item on selection.
     */
    function handleSelect(element: HTMLElement | null): void {
        if (element) {
            const anchor =
                element.tagName.toLowerCase() === 'a'
                    ? (element as HTMLAnchorElement)
                    : element.querySelector<HTMLAnchorElement>('a')
            if (anchor) {
                if (anchor.dataset.goUp) {
                    let currentTreeProvider = treeProvider
                    treeProvider.fetchParent().then(parentTreeProvider => {
                        if (treeProvider === currentTreeProvider) {
                            treeProvider = parentTreeProvider
                        }
                    })
                } else {
                    goto(anchor.href, { keepFocus: true })
                }
            }
        }
    }

    function handleScopeChange(scopedTreeProvider: FileTreeProvider): void {
        treeProvider = scopedTreeProvider.copy({ parent: undefined })
        const root = treeProvider.getRoot()

        if (root === NODE_LIMIT) {
            return
        }

        if (!selectedPath.startsWith(root.path)) {
            goto(replaceRevisionInURL(root.canonicalURL, revision), { keepFocus: true })
        }
    }

    /**
     * For a given path (e.g. foo/bar/baz) returns a list of ancestor paths (e.g.
     * [foo, foo/bar]
     */
    function getAncestorPaths(path: string) {
        return path
            .split('/')
            .slice(0, -1)
            .map((_, index, segements) => segements.slice(0, index + 1).join('/'))
    }

    /**
     * Takes a file path and makes all intermediate nodes as open, and the last node as selected.
     */
    async function markSelected(path: string) {
        const nodesCopy = new Set($treeState.expandedNodes)

        for (const ancestor of getAncestorPaths(path)) {
            nodesCopy.add(ancestor)
        }
        nodesCopy.add(path)

        $treeState = { focused: path, selected: path, expandedNodes: nodesCopy }
    }

    // Since context is only set once when the component is created
    // we need to dynamically sync any changes to the corresponding
    // file tree state store
    const treeState = createForwardStore(getSidebarFileTreeStateForRepo(repoName))

    // Propagating the tree state via context yielded better performance than passing
    // it via props.
    setTreeContext(treeState)

    $: treeRoot = treeProvider.getRoot()
    $: treeState.updateStore(getSidebarFileTreeStateForRepo(repoName))

    // Update open and selected nodes when the path changes.
    $: markSelected(selectedPath)
</script>

<div tabindex="-1">
    <TreeView
        {treeProvider}
        on:select={event => handleSelect(event.detail)}
        on:scope-change={event => handleScopeChange(event.detail.provider)}
    >
        <svelte:fragment let:entry let:expanded>
            {@const isRoot = entry === treeRoot}
            {#if entry === NODE_LIMIT}
                <!-- todo: create alert component -->
                <span class="note">Full list is too long to display. Use search to find a specific file.</span>
            {:else}
                <!--
                    We handle navigation via the TreeView's select event, to preserve the focus state.
                    Using a link here allows us to benefit from data preloading.
                -->
                <Popover let:registerTrigger placement="right-end" showOnHover>
                    <a
                        href={replaceRevisionInURL(entry.canonicalURL, revision)}
                        on:click|preventDefault={() => {}}
                        tabindex={-1}
                        data-go-up={isRoot ? true : undefined}
                        use:registerTrigger
                        on:mouseover={/* Preload */ () =>
                            fetchPopoverData({ repoName, revision, filePath: entry.path })}
                    >
                        {#if entry.isDirectory}
                            <Icon svgPath={getDirectoryIconPath(entry, expanded)} inline />
                        {:else}
                            <FileIcon inline file={entry.__typename === 'GitBlob' ? entry : null} />
                        {/if}
                        {isRoot ? '..' : entry.name}
                    </a>
                    <svelte:fragment slot="content">
                        {#await fetchPopoverData({ repoName, revision, filePath: entry.path }) then entry}
                            <FilePopover {repoName} {revision} {entry} />
                        {/await}
                    </svelte:fragment>
                </Popover>
            {/if}
        </svelte:fragment>
        <Alert slot="error" let:error variant="danger">
            Unable to fetch file tree data: {error.message}
        </Alert>
    </TreeView>
</div>

<style lang="scss">
    div {
        overflow: auto;

        :global([data-treeitem][aria-selected]) > :global([data-treeitem-label]) {
            cursor: pointer;

            &:hover {
                background-color: var(--color-bg-3);
            }
        }

        :global([data-treeitem][aria-selected='true']) > :global([data-treeitem-label]) {
            --tree-node-expand-icon-color: var(--body-bg);
            background-color: var(--primary);

            &:hover {
                background-color: var(--primary);
            }
        }
    }

    a {
        text-overflow: ellipsis;
        overflow: hidden;
        white-space: nowrap;
        color: inherit;
        text-decoration: none;
        width: 100%;

        &:hover {
            color: inherit;
            text-decoration: none;
        }
    }

    .note {
        color: var(--body-color);
        overflow: hidden;
        border: 1px solid var(--border-color);
        background-color: var(--subtle-bg);
        padding: 0.25rem;
        border-radius: var(--border-radius);
    }
</style>
