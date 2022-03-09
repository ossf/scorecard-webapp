<template>
  <footer class="footer py-32 md:py-64 bg-pastel-white">
    <div class="flex flex-wrap justify-between items-start container">
      <div class="w-full md:w-1/2">
        <a href="/" title="Home link" class="logo"
          ><Logo
            :class="[
              $route.params.slug === 'trust-security' ? 'logo-white' : 'logo-black',
            ]"
        /></a>
      </div>
      <p class="text-12 text-gray-dark mt-32">Copyright © 2022</p>
    </div>
  </footer>
</template>
<script lang="js" src="./Footer.js" />
