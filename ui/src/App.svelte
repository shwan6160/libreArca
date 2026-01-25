<script lang="ts">
  import { Router, Route, Link } from "svelte-routing";
  import Home from './pages/Home.svelte'
  import WikiPage from './pages/WikiPage.svelte'

  export let url: string = '';

  const routes = {
    '/': Home,
    '/w/:slug': WikiPage,
  }
</script>

<header>
  <h1>libreArca</h1>
  <nav class="nav">
    <a href="/">Home</a>
    <a href="/w/FrontPage">FrontPage</a>
  </nav>
</header>

<main class="main-content">
  <Router {url}>
    <Route path="/" component={Home} />
    <Route path="/w/:slug" let:params>
      <WikiPage slug={params.slug} />
    </Route>
  </Router>
</main>

<div class="left-sidebar"></div>

<div class="right-sidebar"></div>

<footer>
  <p>© 2024 libreArca</p>
</footer>
