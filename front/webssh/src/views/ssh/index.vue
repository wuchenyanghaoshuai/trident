<template>
  <div class="app-container">
    <!-- 使用 Element UI 的表单组件创建一个带有标签和输入框的表单 -->
    <el-form ref="form" :model="form" :inline="true" label-width="120px">
      <el-form-item label="用户名"> <!-- namespace 输入框 -->
        <el-input v-model="form.username" />
      </el-form-item>
      <el-form-item label="IP"> <!-- pod 名称输入框 -->
        <el-input v-model="form.ip" />
      </el-form-item>
      <el-form-item label="端口"> <!-- pod 名称输入框 -->
        <el-input v-model="form.port" />
      </el-form-item>
      <el-form-item label="Command"> <!-- 命令选择框 -->
        <el-select v-model="form.command" placeholder="bash">
          <el-option label="bash" value="bash" />
          <el-option label="sh" value="sh" />
        </el-select>
      </el-form-item>
      <el-form-item> <!-- 提交按钮 -->
        <el-button type="primary" @click="onSubmit">SSH</el-button>
      </el-form-item>
      <div id="terminal" /> <!-- 终端视图容器 -->
    </el-form>
  </div>
</template>

<script>
import { Terminal } from 'xterm' // 导入 xterm 包，用于创建和操作终端对象
import { common as xtermTheme } from 'xterm-style' // 导入 xterm 样式主题
import 'xterm/css/xterm.css' // 导入 xterm CSS 样式
import { FitAddon } from 'xterm-addon-fit' // 导入 xterm fit 插件，用于调整终端大小
import { WebLinksAddon } from 'xterm-addon-web-links' // 导入 xterm web-links 插件，可以捕获 URL 并将其转换为可点击链接
import 'xterm/lib/xterm.js' // 导入 xterm 库

export default {
  data() {
    return {
      form: {
        username: 'root', // 默认命名空间为 "default"
        password: '123', // 默认 shell 命令为 "bash"
        command: 'bash', // 默认 shell 命令为 "bash"
        auth_type: 'pwd', // 默认容器名称为 "nginx"
        ip: '192.168.3.170',
        port: 22
      }
    }
  },
  methods: {
    onSubmit() {
      // 创建一个新的 Terminal 对象
      const xterm = new Terminal({
        theme: xtermTheme,
        rendererType: 'canvas',
        convertEol: true,
        cursorBlink: true
      })

      // 创建并加载 FitAddon 和 WebLinksAddon
      const fitAddon = new FitAddon()
      xterm.loadAddon(fitAddon)
      xterm.loadAddon(new WebLinksAddon())

      // 打开这个终端，并附加到 HTML 元素上
      xterm.open(document.getElementById('terminal'))

      // 调整终端的大小以适应其父元素
      fitAddon.fit()

      // 创建一个新的 WebSocket 连接，并通过 URL 参数传递 pod, namespace, container 和 command 信息
      const ws = new WebSocket('ws://127.0.0.1:8888/ssh?username=' + this.form.username + '&ip=' + this.form.ip + '&port=' + this.form.port + '&command=' + this.form.command)

      // 当 WebSocket 连接打开时，发送一个 resize 消息给服务器，告诉它终端的尺寸
      ws.onopen = function() {
        ws.send(JSON.stringify({
          type: 'resize',
          rows: xterm.rows,
          cols: xterm.cols
        }))
      }

      // 当从服务器收到消息时，写入终端显示
      ws.onmessage = function(evt) {
        xterm.write(evt.data)
      }

      // 当发生错误时，也写入终端显示
      ws.onerror = function(evt) {
        xterm.write(evt.data)
      }

      // 当窗口尺寸变化时，重新调整终端的尺寸，并发送一个新的 resize 消息给服务器
      window.addEventListener('resize', function() {
        fitAddon.fit()
        ws.send(JSON.stringify({
          type: 'resize',
          rows: xterm.rows,
          cols: xterm.cols
        }))
      })

      // 当在终端中键入字符时，发送一个 input 消息给服务器
      xterm.onData((b) => {
        ws.send(JSON.stringify({
          type: 'input',
          text: b
        }))
      })
    }
  }
}
</script>

<style scoped>
.line{
  text-align: center;
}
</style>
