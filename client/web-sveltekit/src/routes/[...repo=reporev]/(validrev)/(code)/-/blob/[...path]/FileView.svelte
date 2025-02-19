<script lang="ts" context="module">
    const graphQLClient = getGraphQLClient()
</script>

<script lang="ts">
    import { mdiClose, mdiFileEyeOutline, mdiMapSearch, mdiWrap, mdiWrapDisabled } from '@mdi/js'
    import { capitalize } from 'lodash'
    import { from } from 'rxjs'
    import { writable } from 'svelte/store'

    import { noOpTelemetryRecorder } from '@sourcegraph/shared/src/telemetry'

    import { goto, preloadData, afterNavigate } from '$app/navigation'
    import { page } from '$app/stores'
    import type { ScrollSnapshot } from '$lib/codemirror/utils'
    import CodeMirrorBlob from '$lib/CodeMirrorBlob.svelte'
    import { isErrorLike, pluralize, SourcegraphURL, type LineOrPositionOrRange } from '$lib/common'
    import { getGraphQLClient, toGraphQLResult } from '$lib/graphql'
    import Icon from '$lib/Icon.svelte'
    import FileHeader from '$lib/repo/FileHeader.svelte'
    import FileIcon from '$lib/repo/FileIcon.svelte'
    import { renderMermaid } from '$lib/repo/mermaid'
    import OpenInEditor from '$lib/repo/open-in-editor/OpenInEditor.svelte'
    import Permalink from '$lib/repo/Permalink.svelte'
    import { createCodeIntelAPI } from '$lib/shared'
    import { isLightTheme, settings } from '$lib/stores'
    import { TELEMETRY_V2_RECORDER } from '$lib/telemetry2'
    import { codeCopiedEvent, SVELTE_LOGGER, SVELTE_TELEMETRY_EVENTS } from '$lib/telemetry'
    import { createPromiseStore, formatBytes } from '$lib/utils'
    import { Alert, Badge, MenuButton, MenuLink } from '$lib/wildcard'
    import markdownStyles from '$lib/wildcard/Markdown.module.scss'

    import type { PageData } from './$types'
    import { FileViewGitBlob, FileViewHighlightedFile } from './FileView.gql'
    import FileViewModeSwitcher from './FileViewModeSwitcher.svelte'
    import OpenInCodeHostAction from './OpenInCodeHostAction.svelte'
    import { CodeViewMode, toCodeViewMode } from './util'

    export let data: Extract<PageData, { type: 'FileView' }>
    export let embedded: boolean = false
    export let disableCodeIntel: boolean = false

    export function capture(): ScrollSnapshot | null {
        return cmblob?.getScrollSnapshot() ?? null
    }

    export function restore(data: ScrollSnapshot | null) {
        initialScrollPosition = data
    }

    const lineWrap = writable<boolean>(false)
    const blobLoader = createPromiseStore<Awaited<PageData['blob']>>()
    const highlightsLoader = createPromiseStore<Awaited<PageData['highlights']>>()

    let blob: FileViewGitBlob | null = null
    let highlights: FileViewHighlightedFile | null = null
    let cmblob: CodeMirrorBlob | null = null
    let initialScrollPosition: ScrollSnapshot | null = null
    let selectedPosition: LineOrPositionOrRange | null = null

    $: ({
        repoName,
        filePath,
        repoURL,
        blameData,
        revision,
        resolvedRevision: { commitID },
        revisionOverride,
    } = data)
    $: blobLoader.set(data.blob)
    $: highlightsLoader.set(data.highlights)

    $: if (!$blobLoader.pending) {
        // Only update highlights and position after the file content has been loaded.
        // While the file content is loading we show the previous file content.
        blob = $blobLoader.value ?? null
        highlights = $highlightsLoader.pending ? null : $highlightsLoader.value ?? null
        selectedPosition = data.lineOrPosition
    }
    $: fileLoadingError = !$blobLoader.pending && $blobLoader.error
    $: fileNotFound = !$blobLoader.pending && !blob

    $: fileViewModeFromURL = toCodeViewMode($page.url.searchParams.get('view'))
    $: isRichFile = !!blob?.richHTML
    $: isBinaryFile = blob?.binary ?? false

    $: showFileModeSwitcher = blob && !isBinaryFile && !embedded
    $: showFormattedView = isRichFile && fileViewModeFromURL === CodeViewMode.Default
    $: showBlameView = fileViewModeFromURL === CodeViewMode.Blame

    $: codeIntelAPI =
        !isBinaryFile && !showFormattedView && !disableCodeIntel
            ? createCodeIntelAPI({
                  settings: setting => (isErrorLike($settings?.final) ? undefined : $settings?.final?.[setting]),
                  requestGraphQL(options) {
                      return from(graphQLClient.query(options.request, options.variables).then(toGraphQLResult))
                  },
                  telemetryRecorder: noOpTelemetryRecorder,
              })
            : null

    function viewModeURL(viewMode: CodeViewMode) {
        switch (viewMode) {
            case CodeViewMode.Code: {
                const url = SourcegraphURL.from($page.url)
                if (isRichFile) {
                    url.setSearchParameter('view', 'raw')
                } else {
                    url.deleteSearchParameter('view')
                }
                return url.toString()
            }
            case CodeViewMode.Blame:
                const url = SourcegraphURL.from($page.url)
                url.setSearchParameter('view', 'blame')
                return url.toString()
            case CodeViewMode.Default:
                return SourcegraphURL.from($page.url).deleteSearchParameter('view').toString()
        }
    }

    function handleCopy(): void {
        SVELTE_LOGGER.log(...codeCopiedEvent('blob-view'))
    }

    function onViewModeChange(event: CustomEvent<CodeViewMode>): void {
        // TODO: track other blob mode
        if (event.detail === CodeViewMode.Blame) {
            SVELTE_LOGGER.log(SVELTE_TELEMETRY_EVENTS.GitBlameEnabled)
            TELEMETRY_V2_RECORDER.recordEvent('repo.gitBlame', 'enable')
        }

        goto(viewModeURL(event.detail), { replaceState: true, keepFocus: true })
    }

    afterNavigate(event => {
        // Only restore scroll position when the user used the browser history to navigate back
        // and forth. When the user reloads the page, in which case SvelteKit will also call
        // Snapshot.restore, we don't want to restore the scroll position. Instead we want the
        // selected line (if any) to scroll into view.
        if (event.type !== 'popstate') {
            initialScrollPosition = null
        }
    })
</script>

{#if embedded}
    <FileHeader type="blob" repoName={data.repoName} path={data.filePath} {revision}>
        <FileIcon slot="icon" file={blob} inline />
        <svelte:fragment slot="actions">
            <slot name="actions" />
        </svelte:fragment>
    </FileHeader>
{:else if revisionOverride}
    <FileHeader type="blob" repoName={data.repoName} path={data.filePath} {revision}>
        <FileIcon slot="icon" file={blob} inline />
    </FileHeader>
{:else}
    <FileHeader type="blob" repoName={data.repoName} path={data.filePath} {revision}>
        <FileIcon slot="icon" file={blob} inline />
        <svelte:fragment slot="actions">
            {#await data.externalServiceType then externalServiceType}
                {#if externalServiceType && !isBinaryFile}
                    <OpenInEditor {externalServiceType} updateUserSetting={data.updateUserSetting} />
                {/if}
            {/await}
            {#if blob}
                <OpenInCodeHostAction data={blob} />
            {/if}
            <Permalink {commitID} />
        </svelte:fragment>
        <svelte:fragment slot="actionmenu">
            <MenuLink href="{repoURL}/-/raw/{filePath}" target="_blank">
                <Icon svgPath={mdiFileEyeOutline} inline /> View raw
            </MenuLink>
            <MenuButton
                on:click={() => lineWrap.update(wrap => !wrap)}
                disabled={fileViewModeFromURL === CodeViewMode.Default && isRichFile}
            >
                <Icon svgPath={$lineWrap ? mdiWrap : mdiWrapDisabled} inline />
                {$lineWrap ? 'Disable' : 'Enable'} wrapping long lines
            </MenuButton>
        </svelte:fragment>
    </FileHeader>
{/if}

{#if revisionOverride}
    <div class="revision-info">
        <Badge variant="link">
            <a href={revisionOverride.canonicalURL}>{revisionOverride.abbreviatedOID}</a>
        </Badge>
        <a href={SourcegraphURL.from($page.url).deleteSearchParameter('rev').toString()}>
            <Icon svgPath={mdiClose} inline />
            <span>Close commit</span>
        </a>
    </div>
{:else if showFileModeSwitcher}
    <div class="file-info">
        <FileViewModeSwitcher
            aria-label="View mode"
            value={fileViewModeFromURL}
            options={isRichFile
                ? [CodeViewMode.Default, CodeViewMode.Code, CodeViewMode.Blame]
                : [CodeViewMode.Default, CodeViewMode.Blame]}
            on:preload={event => preloadData(viewModeURL(event.detail))}
            on:change={onViewModeChange}
        >
            <svelte:fragment slot="label" let:value>
                {value === CodeViewMode.Default ? (isRichFile ? 'Formatted' : 'Code') : capitalize(value)}
            </svelte:fragment>
        </FileViewModeSwitcher>
        {#if blob}
            <code>
                {blob.totalLines}
                {pluralize('line', blob.totalLines)} · {formatBytes(blob.byteSize)}
            </code>
        {/if}
    </div>
{/if}

<div
    class="content"
    class:loading={$blobLoader.pending}
    class:center={fileLoadingError || fileNotFound || isBinaryFile}
>
    {#if fileLoadingError}
        <Alert variant="danger">
            Unable to load file data: {fileLoadingError.message}
        </Alert>
    {:else if fileNotFound}
        <div class="circle">
            <Icon svgPath={mdiMapSearch} --icon-size="80px" />
        </div>
        <h2>File not found</h2>
    {:else if isBinaryFile}
        <Alert variant="info">
            This is a binary file and cannot be displayed.
            <br />
            <a href="{repoURL}/-/raw/{filePath}" target="_blank" download>Download file</a>
        </Alert>
    {:else if blob && showFormattedView}
        <!-- key on the HTML content so renderMermaid gets re-run -->
        {#key blob.richHTML}
            <!-- jupyter is a global style -->
            <div
                use:renderMermaid={{ selector: 'pre:has(code.language-mermaid)', isLightTheme: $isLightTheme }}
                class={`rich jupyter ${markdownStyles.markdown}`}
            >
                {@html blob.richHTML}
            </div>
        {/key}
    {:else if blob}
        <!--
            This ensures that a new CodeMirror instance is created when the file changes.
            This makes the CodeMirror behavior more predictable and avoids issues with
            carrying over state from the previous file.
            Specifically this will make it so that the scroll position is reset to
            `initialScrollPosition` when the file changes.
        -->
        {#key blob.canonicalURL}
            <CodeMirrorBlob
                bind:this={cmblob}
                {initialScrollPosition}
                blobInfo={{
                    ...blob,
                    repoName,
                    commitID,
                    revision: revision ?? '',
                    filePath,
                }}
                highlights={highlights?.lsif ?? ''}
                showBlame={showBlameView}
                blameData={$blameData}
                wrapLines={$lineWrap}
                selectedLines={selectedPosition?.line ? selectedPosition : null}
                on:selectline={({ detail: range }) => {
                    goto(
                        SourcegraphURL.from($page.url.searchParams)
                            .setLineRange(range ? { line: range.line, endLine: range.endLine } : null)
                            .deleteSearchParameter('popover').search
                    )
                }}
                {codeIntelAPI}
                onCopy={handleCopy}
            />
        {/key}
    {/if}
</div>

<style lang="scss">
    .content {
        overflow: auto;
        flex: 1;
        background-color: var(--code-bg);

        &.center {
            display: flex;
            flex-direction: column;
            align-items: center;
        }
    }

    .revision-info,
    .file-info {
        display: flex;
        align-items: baseline;
        gap: 1rem;
        padding: 0.5rem;
        color: var(--text-muted);
        background-color: var(--color-bg-1);
        box-shadow: var(--fileheader-shadow);

        // Allows for its shadow to cascade over the code panel
        z-index: 1;
    }

    .revision-info {
        justify-content: space-between;
        // Increasing the padding makes the switch between the file view and the diff view
        // less jarring (the code view switcher increases the height of the info bar).
        padding: 0.75rem;

        // This is used to avoid having the whitespace being underlined on hover
        a {
            text-decoration: none;

            &:hover span {
                text-decoration: underline;
            }
        }
    }

    .loading {
        filter: blur(1px);
    }

    .rich {
        padding: 1rem;
        max-width: 50rem;
    }

    .circle {
        background-color: var(--color-bg-2);
        border-radius: 50%;
        padding: 1.5rem;
        margin: 1rem;
    }

    .actions {
        margin-left: auto;
    }
</style>
