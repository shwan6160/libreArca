<script lang="ts">
  import { Router, Route, Link } from "svelte-routing";
  import Home from './pages/Home.svelte'
  import Wiki from './pages/Wiki.svelte'
  import Bbs from './pages/Bbs.svelte'
  import Blog from './pages/Blog.svelte'

  export let url: string = '';

  const routes = {
    '/': Home,
    '/w/:slug': Wiki,
    '/b/*': Bbs,
    '/g/:slug': Blog
  }
</script>

<svelte:head>
  <title>{window.__WIKI_CONFIG__.WikiName || "libreArca"}</title>
</svelte:head>

<header>
  <h1>libreArca</h1>
  <nav class="nav">
    <a href="/">Home</a>
    <a href="/w/FrontPage">FrontPage</a>
    <a href="/b">Board</a>
    <a href="/g/FirstBlogPost">Blog</a>
  </nav>
</header>

<main class="main-content">
  <Router {url}>
    <Route path="/" component={Home} />
    <Route path="/w/:slug" let:params>
      <Wiki slug={params.slug} />
    </Route>
    <Route path="/b/*" let:params>
      <Bbs path={params['*']} />
    </Route>
    <Route path="/g/:slug" let:params>
      <Blog slug={params.slug} />
    </Route>
  </Router>
</main>

<div class="left-sidebar"></div>

<div class="right-sidebar"></div>

<footer>
  <p>© 2024 libreArca</p>
</footer>
