<template>
  <footer class="footer py-32 md:py-64 bg-pastel-orange">
    <div class="flex flex-wrap justify-between items-start container">
      <div class="w-full md:w-1/2">
        Logo
        <Navigation nav-type="footer" :nav-list="navigation" />
      </div>
      <div class="w-full md:w-1/2 md:mt-0 mt-32">
        <Navigation nav-type="social" :nav-list="socialLinks" />
      </div>
      <p class="text-12 text-gray-dark mt-32">Copyright © 2022</p>
    </div>
  </footer>
</template>
<script lang="js" src="./Footer.js" />
