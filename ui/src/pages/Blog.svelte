<script lang="ts">
  export let slug = "";
</script>

<svelte:head>
  <title>{slug} - {window.__WIKI_CONFIG__.BlogName || "libreArca"}</title>
</svelte:head>

<h1>{slug}</h1>
<p>여기서 API를 호출해서 내용을 가져오면 됩니다.</p>