import { createApp } from "vue";
import Toast, { PluginOptions } from "vue-toastification";
import App from "./App.vue";
import "vue-toastification/dist/index.css";

const app = createApp(App);

app.use(Toast);

app.mount("#app");