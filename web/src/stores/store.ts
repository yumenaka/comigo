import { createPinia } from "pinia";
import { App } from "vue";

const store = createPinia();

export function setupStore(app: App<Element>) {
  app.use(store);
}

export { store };