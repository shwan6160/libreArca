<script lang="ts">
  import Pages from './bbs/Pages.svelte';
  import Page from './bbs/Page.svelte';

  export let path = "";

  let slug = "";
  let pageId = null;
  let viewComponent = null;

  $: {
    const segments = path.split('/').filter(Boolean);

    if (segments.length === 1) {
      slug = segments[0];
      pageId = null;
      viewComponent = Pages;
    } else if (segments.length >= 2) {
      slug = segments[0];
      pageId = segments[1];
      viewComponent = Page;
    } else {
        // 예외 처리
    }
  }
</script>

<svelte:head>
  <title>{window.__WIKI_CONFIG__.BbsName || "libreArca"}</title>
</svelte:head>

<svelte:component this={viewComponent} {slug} {pageId} />
