<template>
  <div ref="container" class="art-player"></div>
</template>

<script setup>
import Artplayer from 'artplayer'
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'

const props = defineProps({
  url: {
    type: String,
    required: true,
  },
  danmakus: {
    type: Array,
    default: () => [],
  },
})

const container = ref(null)
let player = null

function mountPlayer() {
  if (!container.value) return

  player?.destroy(false)
  player = new Artplayer({
    container: container.value,
    url: props.url,
    autoplay: false,
    autoSize: true,
    fullscreen: true,
    fullscreenWeb: true,
    setting: true,
    playbackRate: true,
    aspectRatio: true,
    miniProgressBar: true,
    mutex: true,
    theme: '#fb7299',
    subtitleOffset: true,
    controls: [
      {
        position: 'right',
        html: 'GoDan',
        tooltip: 'phase 12',
        style: {
          color: '#fff',
          fontWeight: 700,
        },
      },
    ],
  })
}

onMounted(mountPlayer)
watch(() => props.url, mountPlayer)

onBeforeUnmount(() => {
  player?.destroy(false)
  player = null
})
</script>

<style scoped>
.art-player {
  width: 100%;
  aspect-ratio: 16 / 9;
}
</style>
