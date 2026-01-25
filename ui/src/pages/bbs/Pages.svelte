<script lang="ts">
  export let slug: string = "";
</script>

<svelte:head>
  <title>{slug} - {window.__WIKI_CONFIG__.BbsName || "libreArca"}</title>
</svelte:head>

<div class="pages">
    <h3>Pages</h3>
    <ul>
        <li>test 1</li>
        <li>test 2</li>
        <li>test 3</li>
    </ul>
</div>