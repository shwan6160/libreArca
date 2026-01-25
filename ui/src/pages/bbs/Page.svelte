<script lang="ts">
  import Comments from './Comments.svelte';

  export let slug = "";
  export let pageId = null;
</script>

<svelte:head>
  <title>{pageId}제목 - {slug}</title>
</svelte:head>

<h1>{slug}</h1>
<h1>{pageId}의 제목</h1>
<p>여기서 API를 호출해서 내용을 가져오면 됩니다.</p>
<Comments />