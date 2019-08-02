import Vue from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';
import './registerServiceWorker';
import './plugins/element.js'
import setPrototypeOf = Reflect.setPrototypeOf;

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  data: () => {
    return{
      Dots: [
        {
          "Meta": {
            "typeId": "4b8b1751-4799-4578-af46-d9b339cf582f",
            "version": "",
            "name": "",
            "showName": "",
            "single": false,
            "relyTypeIds": null
          },
          "Lives": [
            {
              "TypeId": "",
              "LiveId": "",
              "RelyLives": null,
              "Dot": null,
              "json": null,
              "name": ""
            }
          ]
        },
        {
          "Meta": {
            "typeId": "4943e959-7ad7-42c6-84dd-8b24e9ed30bb",
            "version": "",
            "name": "",
            "showName": "",
            "single": false,
            "relyTypeIds": null
          },
          "Lives": [
            {
              "TypeId": "",
              "LiveId": "",
              "RelyLives": null,
              "Dot": null,
              "json": {
                "addr": "",
                "keyFile": "",
                "pemFile": "",
                "logSkipPaths": null
              },
              "name": ""
            }
          ]
        },
        {
          "Meta": {
            "typeId": "6be39d0b-3f5b-47b4-818c-642c049f3166",
            "version": "",
            "name": "",
            "showName": "",
            "single": false,
            "relyTypeIds": null
          },
          "Lives": [
            {
              "TypeId": "",
              "LiveId": "",
              "RelyLives": null,
              "Dot": null,
              "json": {
                "relativePath": ""
              },
              "name": ""
            }
          ]
        },
        {
          "Meta": {
            "typeId": "2d281e77-3133-4718-a5b7-9eea069591ca",
            "version": "",
            "name": "",
            "showName": "",
            "single": false,
            "relyTypeIds": null
          },
          "Lives": [
            {
              "TypeId": "",
              "LiveId": "",
              "RelyLives": null,
              "Dot": null,
              "json": {
                "name": ""
              },
              "name": ""
            }
          ]
        },
        {
          "Meta": {
            "typeId": "7bf0a017-ef0c-496a-b04c-b1dc262abc8d",
            "version": "",
            "name": "",
            "showName": "",
            "single": false,
            "relyTypeIds": null
          },
          "Lives": [
            {
              "TypeId": "",
              "LiveId": "",
              "RelyLives": null,
              "Dot": null,
              "json": {
                "scheme": "",
                "services": [
                  {
                    "name": "",
                    "addrs": null,
                    "tls": {
                      "caPem": "",
                      "pem": "",
                      "key": "",
                      "serverNameOverride": ""
                    },
                    "balance": ""
                  }
                ]
              },
              "name": ""
            }
          ]
        },
        {
          "Meta": {
            "typeId": "3c9e8119-3d42-45bd-98f9-32939c895c6d",
            "version": "",
            "name": "",
            "showName": "",
            "single": false,
            "relyTypeIds": null
          },
          "Lives": [
            {
              "TypeId": "",
              "LiveId": "3c9e8119-3d42-45bd-98f9-32939c895c6d",
              "RelyLives": {
                "GinRouter": "6be39d0b-3f5b-47b4-818c-642c049f3166",
                "ServerNobl": "77a766e7-c288-413f-946b-bc9de6df3d70"
              },
              "Dot": null,
              "json": null,
              "name": ""
            }
          ]
        },
        {
          "Meta": {
            "typeId": "77a766e7-c288-413f-946b-bc9de6df3d70",
            "version": "",
            "name": "",
            "showName": "",
            "single": false,
            "relyTypeIds": null
          },
          "Lives": [
            {
              "TypeId": "",
              "LiveId": "",
              "RelyLives": null,
              "Dot": null,
              "json": {
                "addrs": null,
                "tls": {
                  "caPem": "",
                  "pem": "",
                  "key": "",
                  "serverNameOverride": ""
                }
              },
              "name": ""
            }
          ]
        },
        {
          "Meta": {
            "typeId": "afbeac47-e5fd-4bf3-8fb1-f0fb8ec79bd0",
            "version": "",
            "name": "",
            "showName": "",
            "single": false,
            "relyTypeIds": null
          },
          "Lives": [
            {
              "TypeId": "",
              "LiveId": "afbeac47-e5fd-4bf3-8fb1-f0fb8ec79bd0",
              "RelyLives": {
                "ServerNobl": "77a766e7-c288-413f-946b-bc9de6df3d70"
              },
              "Dot": null,
              "json": {
                "preUrl": "",
                "addr": "",
                "tls": {
                  "caPem": "",
                  "pem": "",
                  "key": "",
                  "serverNameOverride": ""
                }
              },
              "name": ""
            }
          ]
        }
      ],
      Configs: [
        {
          "Meta": {
            "typeId": "4b8b1751-4799-4578-af46-d9b339cf582f",
            "version": "",
            "name": "",
            "showName": "",
            "single": false,
            "relyTypeIds": null
          },
          "Lives": [
            {
              "TypeId": "",
              "LiveId": "",
              "RelyLives": null,
              "Dot": null,
              "json": null,
              "name": ""
            }
          ]
        },
        {
          "Meta": {
            "typeId": "4943e959-7ad7-42c6-84dd-8b24e9ed30bb",
            "version": "",
            "name": "",
            "showName": "",
            "single": false,
            "relyTypeIds": null
          },
          "Lives": [
            {
              "TypeId": "",
              "LiveId": "",
              "RelyLives": null,
              "Dot": null,
              "json": {
                "addr": "",
                "keyFile": "",
                "pemFile": "",
                "logSkipPaths": null
              },
              "name": ""
            }
          ]
        },
        {
          "Meta": {
            "typeId": "6be39d0b-3f5b-47b4-818c-642c049f3166",
            "version": "",
            "name": "",
            "showName": "",
            "single": false,
            "relyTypeIds": null
          },
          "Lives": [
            {
              "TypeId": "",
              "LiveId": "",
              "RelyLives": null,
              "Dot": null,
              "json": {
                "relativePath": ""
              },
              "name": ""
            }
          ]
        },
        {
          "Meta": {
            "typeId": "2d281e77-3133-4718-a5b7-9eea069591ca",
            "version": "",
            "name": "",
            "showName": "",
            "single": false,
            "relyTypeIds": null
          },
          "Lives": [
            {
              "TypeId": "",
              "LiveId": "",
              "RelyLives": null,
              "Dot": null,
              "json": {
                "name": ""
              },
              "name": ""
            }
          ]
        },
        {
          "Meta": {
            "typeId": "7bf0a017-ef0c-496a-b04c-b1dc262abc8d",
            "version": "",
            "name": "",
            "showName": "",
            "single": false,
            "relyTypeIds": null
          },
          "Lives": [
            {
              "TypeId": "",
              "LiveId": "",
              "RelyLives": null,
              "Dot": null,
              "json": {
                "scheme": "",
                "services": [
                  {
                    "name": "",
                    "addrs": null,
                    "tls": {
                      "caPem": "",
                      "pem": "",
                      "key": "",
                      "serverNameOverride": ""
                    },
                    "balance": ""
                  }
                ]
              },
              "name": ""
            }
          ]
        },
        {
          "Meta": {
            "typeId": "3c9e8119-3d42-45bd-98f9-32939c895c6d",
            "version": "",
            "name": "",
            "showName": "",
            "single": false,
            "relyTypeIds": null
          },
          "Lives": [
            {
              "TypeId": "",
              "LiveId": "3c9e8119-3d42-45bd-98f9-32939c895c6d",
              "RelyLives": {
                "GinRouter": "6be39d0b-3f5b-47b4-818c-642c049f3166",
                "ServerNobl": "77a766e7-c288-413f-946b-bc9de6df3d70"
              },
              "Dot": null,
              "json": null,
              "name": ""
            }
          ]
        },
        {
          "Meta": {
            "typeId": "77a766e7-c288-413f-946b-bc9de6df3d70",
            "version": "",
            "name": "",
            "showName": "",
            "single": false,
            "relyTypeIds": null
          },
          "Lives": [
            {
              "TypeId": "",
              "LiveId": "",
              "RelyLives": null,
              "Dot": null,
              "json": {
                "addrs": null,
                "tls": {
                  "caPem": "",
                  "pem": "",
                  "key": "",
                  "serverNameOverride": ""
                }
              },
              "name": ""
            }
          ]
        },
        {
          "Meta": {
            "typeId": "afbeac47-e5fd-4bf3-8fb1-f0fb8ec79bd0",
            "version": "",
            "name": "",
            "showName": "",
            "single": false,
            "relyTypeIds": null
          },
          "Lives": [
            {
              "TypeId": "",
              "LiveId": "afbeac47-e5fd-4bf3-8fb1-f0fb8ec79bd0",
              "RelyLives": {
                "ServerNobl": "77a766e7-c288-413f-946b-bc9de6df3d70"
              },
              "Dot": null,
              "json": {
                "preUrl": "",
                "addr": "",
                "tls": {
                  "caPem": "",
                  "pem": "",
                  "key": "",
                  "serverNameOverride": ""
                }
              },
              "name": ""
            }
          ]
        }
      ]
    }
  },
  render: (h) => h(App),
}).$mount('#app');
