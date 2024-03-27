<template>
  <div id="app">
    <div v-if="!editMode">
      <div>监听地址: {{ listenHost }}</div>
      <div>监听端口: {{ listenPort }}</div>
      <div>订阅链接: {{ subscriptionUrl }}</div>
      <div>节点过滤: {{ keyword }}</div>
      <Button @click="enableEditMode">修改</Button>
      <Button @click="updateSub">更新节点</Button>
    </div>
    <div v-else>
      <Input v-model="form.host" />
      <Input v-model.number="form.port" type="number" />
      <Input v-model="form.subscriptionUrl" />
      <Input v-model="form.keyword" />
      <Button v-if="editMode" @click="updateConfig()">保存</Button>
      <Button v-if="editMode" @click="editMode = false">取消</Button>
    </div>

    <Checkbox v-model="autoMode" @change="changeMode">自动模式</Checkbox>
    <Button @click="next" :disabled="!autoMode">换一个节点</Button>
    <RadioGroup :style="{ display: 'block' }" v-model="radio" v-for="(param, index) in params">
      <Radio :disabled="autoMode" :label="index" :value="param.remarks" @change="picked(param.remarks)">{{ param.remarks }} TTL: {{ param.ttl }}</Radio>
    </RadioGroup>
  </div>
</template>

<script>
import axios from 'axios'
import { Radio, RadioGroup, Input, Button, Switch, Checkbox, Dialog } from 'element-ui'
import { log } from './utils'

import './app.css'

export default {
  components: { Radio, RadioGroup, Input, Button, Switch, Checkbox, Dialog },

  data() {
    return {
      listenHost: '',
      listenPort: 1080,
      subscriptionUrl: '',
      keyword: '',
      autoMode: true,
      radio: 0,
      currentParam: undefined,
      params: [],
      editMode: false,
      form: {
        host: '',
        port: 1080,
        subscriptionUrl: '',
        keyword: '',
      }
    }
  },

  async beforeMount() {
    await this.updateStatus()
    setInterval(this.updateStatus, 10000)
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
        auto_mode: autoMode,
      } = r.data
      let goodParams = params.filter(p => p.ttl > 0)
      goodParams.sort((a, b) => a.ttl - b.ttl)
      let banParams = params.filter(p => p.ttl < 1)

      this.params = [...goodParams, ...banParams]
      let paramRemarks = this.params.map(p => p.remarks)
      if (paramRemarks.length > 0) {
        this.radio = paramRemarks.indexOf(currentParam.remarks)
      }
      this.listenHost = listenHost
      this.listenPort = listenPort
      this.subscriptionUrl = subscriptionUrl
      this.keyword = keyword
      this.autoMode = autoMode
    },

    enableEditMode() {
      this.form.host = this.listenHost
      this.form.port = this.listenPort
      this.form.subscriptionUrl = this.subscriptionUrl
      this.form.keyword = this.keyword
      this.editMode = true
    },

    async updateSub() {
      await axios.get('api/subscription/update')
      setTimeout(this.updateStatus, 1000)
    },

    async updateConfig() {
      await axios.post('api/config', this.form)
    },

    async next() {
      await axios.post('/api/next')
      await this.updateStatus()
    },

    async changeMode(auto_mode_flag) {
      await axios.post('/api/mode', {mode: auto_mode_flag ? 'auto': 'static'})
    },

    async picked(paramName) {
      await axios.post('/api/mode', {mode: 'static', remarks: paramName})
    }
  },

}
</script>
