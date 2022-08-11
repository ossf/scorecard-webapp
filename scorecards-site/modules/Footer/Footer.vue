<template>
  <footer class="footer py-32 md:py-64">
    <div
      class="flex flex-col flex-flex-wrap justify-between items-start container"
    >
      <a href="/" title="Home link" class="logo mb-24"
        ><Logo
          :class="[
            $route.params.slug === 'trust-security'
              ? 'logo-white'
              : 'logo-black',
          ]"
      /></a>
      <div class="text-14 mb-12">
        © 2022 The Linux Foundation, under the terms of the Apache License 2.0.
      </div>
      <div class="text-14">
        <strong>Privacy statement: </strong> We use <a href="https://www.goatcounter.com/">GoatCounter</a> to help us improve our website by
        collecting and reporting information on how it's used. We do not store
        advertising or tracking cookies. The information we collect does not
        identify anyone and does not track an individual's use of the site.
      </div>
      <div>
        <br>
      	<a href="https://www.netlify.com"> <img src="https://www.netlify.com/v3/img/components/netlify-dark.svg" alt="Deploys by Netlify"/> </a>
      </div>
    </div>
  </footer>
</template>
<script lang="js" src="./Footer.js" />
