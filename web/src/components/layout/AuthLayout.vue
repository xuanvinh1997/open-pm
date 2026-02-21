<template>
  <div class="flex min-h-screen items-center justify-center bg-custom-background-80 px-4">
    <div class="w-full max-w-[400px]">
      <div class="mb-10 text-center">
        <div class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-brand-600">
          <span class="text-xl font-bold text-white">O</span>
        </div>
        <h1 class="text-2xl font-bold text-custom-text-100">Open-PM</h1>
        <p class="mt-1 text-sm text-custom-text-300">Project management, simplified.</p>
      </div>
      <div class="rounded-xl border border-custom-border-200 bg-custom-background-100 p-8 shadow-custom-sm">
        <router-view />
      </div>
    </div>
  </div>
</template>
