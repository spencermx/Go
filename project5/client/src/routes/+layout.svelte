<script>
  import { page } from '$app/stores';
</script>

<header class="banner">
  <div class="container">
    <h1>My Video Gallery</h1>
    <nav>
      <ul>
        <li><a href="/" class:active={$page.url.pathname === '/'}>Home</a></li>
        <li><a href="/about" class:active={$page.url.pathname === '/about'}>About</a></li>
        <li><a href="/contact" class:active={$page.url.pathname === '/contact'}>Contact</a></li>
      </ul>
    </nav>
  </div>
</header>

<main>
  <slot></slot>
</main>

<style>
  .banner {
    background-color: #333;
    padding: 1rem 0;
  }

  .banner h1 {
    color: white;
    margin: 0;
    font-size: 2rem;
    font-weight: bold;
  }

  .banner nav ul {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    justify-content: flex-end;
  }

  .banner nav li {
    margin-left: 1.5rem;
  }

  .banner nav a {
    color: white;
    text-decoration: none;
    font-size: 1.1rem;
    transition: color 0.3s ease;
  }

  .banner nav a:hover {
    color: #ffcc00;
  }

  .banner nav a.active {
    font-weight: bold;
    border-bottom: 2px solid #ffcc00;
  }

  .footer {
    background-color: #f8f8f8;
    padding: 1rem 0;
    margin-top: 2rem;
    text-align: center;
    color: #666;
    font-size: 0.9rem;
  }

  .container {
    max-width: 960px;
    margin: 0 auto;
    padding: 0 1rem;
  }

  main {
    padding: 2rem 0;
  }
</style>
