<template>
  <div id="app">
    <Input :value="listenHost" />
    <Input :value="listenPort" />
    <Input :value="subscriptionUrl" />
    <Button @click="next">Next</Button>
    <Input :value="keyword" />
    <RadioGroup :style="{ display: 'block' }" v-model="radio" v-for="(param, index) in params">
      <Radio :label="index">{{ param.remarks }} TTL: {{ param.ttl }}</Radio>
    </RadioGroup>
  </div>
</template>

<script>
import axios from 'axios'
import { Radio, RadioGroup, Input, Button } from 'element-ui'
import { log } from './utils'

import './app.css'

export default {
  components: { Radio, RadioGroup, Input, Button },
  data() {
    return {
      listenHost: '',
      listenPort: 1080,
      subscriptionUrl: '',
      keyword: '',
      radio: 0,
      currentParam: undefined,
      params: []
    }
  },
  methods: {
    async updateStatus() {
      log('update status')
      let r = await axios.get('/api/status')
      let {
        listen_host: listenHost,
        listen_port: listenPort,
        subscription_url: subscriptionUrl,
        current_ssr_param: currentParam,
        ssr_params: params,
        keyword: keyword,
      } = r.data
      let goodParams = params.filter(p => p.ttl > 0)
      goodParams.sort((a, b) => a.ttl - b.ttl)
      let paramRemarks = goodParams.map(p => p.remarks)
      let currentIndex = paramRemarks.indexOf(currentParam.remarks)
      let banParams = params.filter(p => p.ttl === -1)

      this.params = [...goodParams, ...banParams]
      this.listenHost = listenHost
      this.listenPort = listenPort
      this.radio = currentIndex
      this.subscriptionUrl = subscriptionUrl
      this.keyword = keyword
    },
    async next() {
      await axios.post('/api/next')
      await this.updateStatus()
    }
  },
  async beforeMount() {
    await this.updateStatus()
    setInterval(this.updateStatus, 5000)
  }
}
</script>
