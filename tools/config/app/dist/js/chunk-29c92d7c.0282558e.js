(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-29c92d7c"],{"0b6e":function(e,t,n){},"0d51":function(e,t,n){"use strict";n.r(t);var r=n("6422"),a=n.n(r);for(var i in r)"default"!==i&&function(e){n.d(t,e,function(){return r[e]})}(i);t["default"]=a.a},1071:function(e,t,n){"use strict";n.r(t);var r=n("c034"),a=n("31e3");for(var i in a)"default"!==i&&function(e){n.d(t,e,function(){return a[e]})}(i);n("d4c1");var o=n("2877"),s=Object(o["a"])(a["default"],r["a"],r["b"],!1,null,"2616f3ce",null);t["default"]=s.exports},"11e9":function(e,t,n){var r=n("52a7"),a=n("4630"),i=n("6821"),o=n("6a99"),s=n("69a8"),l=n("c69a"),c=Object.getOwnPropertyDescriptor;t.f=n("9e1e")?c:function(e,t){if(e=i(e),t=o(t,!0),l)try{return c(e,t)}catch(n){}if(s(e,t))return a(!r.f.call(e,t),e[t])}},2338:function(e,t,n){"use strict";n.r(t);var r=n("bb3c"),a=n("7939");for(var i in a)"default"!==i&&function(e){n.d(t,e,function(){return a[e]})}(i);var o=n("2877"),s=Object(o["a"])(a["default"],r["a"],r["b"],!1,null,"20e8bc5c",null);t["default"]=s.exports},2366:function(e,t){for(var n=[],r=0;r<256;++r)n[r]=(r+256).toString(16).substr(1);function a(e,t){var r=t||0,a=n;return[a[e[r++]],a[e[r++]],a[e[r++]],a[e[r++]],"-",a[e[r++]],a[e[r++]],"-",a[e[r++]],a[e[r++]],"-",a[e[r++]],a[e[r++]],"-",a[e[r++]],a[e[r++]],a[e[r++]],a[e[r++]],a[e[r++]],a[e[r++]]].join("")}e.exports=a},"31e3":function(e,t,n){"use strict";n.r(t);var r=n("3d60"),a=n.n(r);for(var i in r)"default"!==i&&function(e){n.d(t,e,function(){return r[e]})}(i);t["default"]=a.a},"37c8":function(e,t,n){t.f=n("2b4c")},"3a72":function(e,t,n){var r=n("7726"),a=n("8378"),i=n("2d00"),o=n("37c8"),s=n("86cc").f;e.exports=function(e){var t=a.Symbol||(a.Symbol=i?{}:r.Symbol||{});"_"==e.charAt(0)||e in t||s(t,e,{value:o.f(e)})}},"3d60":function(e,t,n){"use strict";var r=n("288e");n("1c01"),Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var a=r(n("7618"));n("ac4d"),n("8a81"),n("ac6a");var i=r(n("2b0e")),o=r(n("2338")),s=r(n("51d2")),l=r(n("99a9")),c=n("c437"),u=i.default.extend({data:function(){return{activeTypes:[],dialogVisible:!1,textarea:"",keyarea:"",objc:null,model:"",schemaObject:{},upLoadFile:"",typeId:"",lives:[]}},components:{"rely-lives-editor":o.default,"json-button":s.default,"check-live-id":l.default},methods:{removeType:function(e){this.$root.Configs.splice(e,1)},AddObject:function(e,t){if(e.metaData.flag){var n=this.shallowCopy(e.lives[e.lives.length-1]);e.lives.push(n)}else{var r=!0,a=!1,i=void 0;try{for(var o,s=this.$root.Dots[Symbol.iterator]();!(r=(o=s.next()).done);r=!0){var l=o.value;if(l.metaData.typeId===t){var c=this.shallowCopy(l.lives[0]);e.lives.push(c);break}}}catch(u){a=!0,i=u}finally{try{r||null==s.return||s.return()}finally{if(a)throw i}}}},RemoveObject:function(e,t){e.lives.splice(t,1)},UuidGenerator:function(e){e.liveId=c()},shallowCopy:function(e){var t;return t=JSON.stringify(e),JSON.parse(t)},uploadSectionFile:function(e){var t=e.file,n=(t.slice(),new FileReader);n.readAsText(t,"utf-8"),n.onload=this.fileOnload},fileOnload:function(e){this.upLoadFile=e.target.result},showDialog:function(e,t){this.typeId=e,this.lives=t,this.dialogVisible=!0},handleConfirm:function(){var e=this.findConfigLives(this.typeId),t=this.findDotLive(this.typeId);if(e&&t)for(var n=0,r=e.length;n<r;n++)this.assembleByLiveId(t,e[n]);this.dialogVisible=!1},findConfigLives:function(e){for(var t=JSON.parse(this.upLoadFile),n=0,r=t.dots.length;n<r;n++)if(t.dots[n].metaData.typeId===e)return t.dots[n].lives;return null},findDotLive:function(e){for(var t=0,n=this.$root.Dots.length;t<n;t++)if(this.$root.Dots[t].metaData.typeId===e)return this.$root.Dots[t].lives[0];return null},assemble:function(e,t){for(var n in e)e[n]=this.isObject(e[n])?this.assemble(e[n],t[n]?t[n]:e[n]):t[n]?t[n]:e[n];return e},isObject:function(e){return"object"===(0,a.default)(e)&&null!==e},assembleByLiveId:function(e,t){for(var n=!1,r=0,a=this.lives.length;r<a;r++)if(this.lives[r].liveId===t.liveId){this.lives[r]=JSON.parse(JSON.stringify(this.assemble(this.lives[r],t))),n=!0;break}if(!n){var i=this.assemble(e,t);this.lives.push(JSON.parse(JSON.stringify(i)))}}}});t.default=u},"51d2":function(e,t,n){"use strict";n.r(t);var r=n("7d38"),a=n("0d51");for(var i in a)"default"!==i&&function(e){n.d(t,e,function(){return a[e]})}(i);var o=n("2877"),s=Object(o["a"])(a["default"],r["a"],r["b"],!1,null,"6ac3897d",null);t["default"]=s.exports},5977:function(e,t,n){"use strict";n.r(t);var r=n("c39b"),a=n.n(r);for(var i in r)"default"!==i&&function(e){n.d(t,e,function(){return r[e]})}(i);t["default"]=a.a},6422:function(e,t,n){"use strict";var r=n("288e");n("1c01"),Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var a=r(n("7618")),i=r(n("2b0e")),o=i.default.extend({name:"JsonButton",props:{objc:{type:Object,required:!0}},data:function(){return{dialog:!1,schemaObject:{},textarea:""}},methods:{ShowJsonDialog:function(e){this.dialog=!0,this.objc=e;var t=n("7320");for(var r in this.schemaObject=t.jsonToSchema(e),this.schemaObject.properties)if("relyLives"===r){this.schemaObject.properties[r].required=[],this.$delete(this.schemaObject.properties[r],"maxProperties");break}this.textarea=JSON.stringify(e,null,4)},handleClose:function(e){try{if(this.textarea){var t=JSON.parse(this.textarea),r=n("8e52");r.validate(t,this.schemaObject)&&this.relyLivesValidation(t)?this.$emit("input",t):this.$message.error("json text input error!")}else this.$message.error("json text input error!")}catch(a){this.$message.error("json text input error!")}finally{e()}},relyLivesValidation:function(e){for(var t in e)if("relyLives"===t)return this.isRelyLives(e[t]);return!0},isRelyLives:function(e){for(var t in e)if((0,a.default)(e[t])!==(0,a.default)(""))return!1;return!0}}});t.default=o},"67ab":function(e,t,n){var r=n("ca5a")("meta"),a=n("d3f4"),i=n("69a8"),o=n("86cc").f,s=0,l=Object.isExtensible||function(){return!0},c=!n("79e5")(function(){return l(Object.preventExtensions({}))}),u=function(e){o(e,r,{value:{i:"O"+ ++s,w:{}}})},f=function(e,t){if(!a(e))return"symbol"==typeof e?e:("string"==typeof e?"S":"P")+e;if(!i(e,r)){if(!l(e))return"F";if(!t)return"E";u(e)}return e[r].i},d=function(e,t){if(!i(e,r)){if(!l(e))return!0;if(!t)return!1;u(e)}return e[r].w},v=function(e){return c&&p.NEED&&l(e)&&!i(e,r)&&u(e),e},p=e.exports={KEY:r,NEED:!1,fastKey:f,getWeak:d,onFreeze:v}},7939:function(e,t,n){"use strict";n.r(t);var r=n("ccc3"),a=n.n(r);for(var i in r)"default"!==i&&function(e){n.d(t,e,function(){return r[e]})}(i);t["default"]=a.a},"7bbc":function(e,t,n){var r=n("6821"),a=n("9093").f,i={}.toString,o="object"==typeof window&&window&&Object.getOwnPropertyNames?Object.getOwnPropertyNames(window):[],s=function(e){try{return a(e)}catch(t){return o.slice()}};e.exports.f=function(e){return o&&"[object Window]"==i.call(e)?s(e):a(r(e))}},"7d38":function(e,t,n){"use strict";var r=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",[n("el-button",{on:{click:function(t){return e.ShowJsonDialog(e.objc)}}},[e._v("JSON")]),n("el-drawer",{ref:"drawer",attrs:{title:"JSON textarea!","before-close":e.handleClose,visible:e.dialog,direction:"ltr","custom-class":"demo-drawer"},on:{"update:visible":function(t){e.dialog=t}}},[n("el-input",{attrs:{type:"textarea",autosize:{minRows:10,maxRows:30},placeholder:"Please input json data!"},model:{value:e.textarea,callback:function(t){e.textarea=t},expression:"textarea"}})],1)],1)},a=[];n.d(t,"a",function(){return r}),n.d(t,"b",function(){return a})},8967:function(e,t,n){"use strict";n.r(t);var r=n("9c5e"),a=n("c9d2");for(var i in a)"default"!==i&&function(e){n.d(t,e,function(){return a[e]})}(i);var o=n("2877"),s=Object(o["a"])(a["default"],r["a"],r["b"],!1,null,"6d5ca63c",null);t["default"]=s.exports},"8a81":function(e,t,n){"use strict";var r=n("7726"),a=n("69a8"),i=n("9e1e"),o=n("5ca1"),s=n("2aba"),l=n("67ab").KEY,c=n("79e5"),u=n("5537"),f=n("7f20"),d=n("ca5a"),v=n("2b4c"),p=n("37c8"),h=n("3a72"),b=n("d4c0"),m=n("1169"),y=n("cb7c"),g=n("d3f4"),j=n("4bf8"),x=n("6821"),w=n("6a99"),O=n("4630"),_=n("2aeb"),k=n("7bbc"),S=n("11e9"),D=n("2621"),C=n("86cc"),I=n("0d58"),N=S.f,$=C.f,J=k.f,L=r.Symbol,P=r.JSON,R=P&&P.stringify,E="prototype",F=v("_hidden"),T=v("toPrimitive"),V={}.propertyIsEnumerable,q=u("symbol-registry"),A=u("symbols"),M=u("op-symbols"),z=Object[E],B="function"==typeof L&&!!D.f,K=r.QObject,G=!K||!K[E]||!K[E].findChild,U=i&&c(function(){return 7!=_($({},"a",{get:function(){return $(this,"a",{value:7}).a}})).a})?function(e,t,n){var r=N(z,t);r&&delete z[t],$(e,t,n),r&&e!==z&&$(z,t,r)}:$,W=function(e){var t=A[e]=_(L[E]);return t._k=e,t},Y=B&&"symbol"==typeof L.iterator?function(e){return"symbol"==typeof e}:function(e){return e instanceof L},Q=function(e,t,n){return e===z&&Q(M,t,n),y(e),t=w(t,!0),y(n),a(A,t)?(n.enumerable?(a(e,F)&&e[F][t]&&(e[F][t]=!1),n=_(n,{enumerable:O(0,!1)})):(a(e,F)||$(e,F,O(1,{})),e[F][t]=!0),U(e,t,n)):$(e,t,n)},H=function(e,t){y(e);var n,r=b(t=x(t)),a=0,i=r.length;while(i>a)Q(e,n=r[a++],t[n]);return e},X=function(e,t){return void 0===t?_(e):H(_(e),t)},Z=function(e){var t=V.call(this,e=w(e,!0));return!(this===z&&a(A,e)&&!a(M,e))&&(!(t||!a(this,e)||!a(A,e)||a(this,F)&&this[F][e])||t)},ee=function(e,t){if(e=x(e),t=w(t,!0),e!==z||!a(A,t)||a(M,t)){var n=N(e,t);return!n||!a(A,t)||a(e,F)&&e[F][t]||(n.enumerable=!0),n}},te=function(e){var t,n=J(x(e)),r=[],i=0;while(n.length>i)a(A,t=n[i++])||t==F||t==l||r.push(t);return r},ne=function(e){var t,n=e===z,r=J(n?M:x(e)),i=[],o=0;while(r.length>o)!a(A,t=r[o++])||n&&!a(z,t)||i.push(A[t]);return i};B||(L=function(){if(this instanceof L)throw TypeError("Symbol is not a constructor!");var e=d(arguments.length>0?arguments[0]:void 0),t=function(n){this===z&&t.call(M,n),a(this,F)&&a(this[F],e)&&(this[F][e]=!1),U(this,e,O(1,n))};return i&&G&&U(z,e,{configurable:!0,set:t}),W(e)},s(L[E],"toString",function(){return this._k}),S.f=ee,C.f=Q,n("9093").f=k.f=te,n("52a7").f=Z,D.f=ne,i&&!n("2d00")&&s(z,"propertyIsEnumerable",Z,!0),p.f=function(e){return W(v(e))}),o(o.G+o.W+o.F*!B,{Symbol:L});for(var re="hasInstance,isConcatSpreadable,iterator,match,replace,search,species,split,toPrimitive,toStringTag,unscopables".split(","),ae=0;re.length>ae;)v(re[ae++]);for(var ie=I(v.store),oe=0;ie.length>oe;)h(ie[oe++]);o(o.S+o.F*!B,"Symbol",{for:function(e){return a(q,e+="")?q[e]:q[e]=L(e)},keyFor:function(e){if(!Y(e))throw TypeError(e+" is not a symbol!");for(var t in q)if(q[t]===e)return t},useSetter:function(){G=!0},useSimple:function(){G=!1}}),o(o.S+o.F*!B,"Object",{create:X,defineProperty:Q,defineProperties:H,getOwnPropertyDescriptor:ee,getOwnPropertyNames:te,getOwnPropertySymbols:ne});var se=c(function(){D.f(1)});o(o.S+o.F*se,"Object",{getOwnPropertySymbols:function(e){return D.f(j(e))}}),P&&o(o.S+o.F*(!B||c(function(){var e=L();return"[null]"!=R([e])||"{}"!=R({a:e})||"{}"!=R(Object(e))})),"JSON",{stringify:function(e){var t,n,r=[e],a=1;while(arguments.length>a)r.push(arguments[a++]);if(n=t=r[1],(g(t)||void 0!==e)&&!Y(e))return m(t)||(t=function(e,t){if("function"==typeof n&&(t=n.call(this,e,t)),!Y(t))return t}),r[1]=t,R.apply(P,r)}}),L[E][T]||n("32e9")(L[E],T,L[E].valueOf),f(L,"Symbol"),f(Math,"Math",!0),f(r.JSON,"JSON",!0)},9093:function(e,t,n){var r=n("ce10"),a=n("e11e").concat("length","prototype");t.f=Object.getOwnPropertyNames||function(e){return r(e,a)}},"99a9":function(e,t,n){"use strict";n.r(t);var r=n("9d7e0"),a=n("5977");for(var i in a)"default"!==i&&function(e){n.d(t,e,function(){return a[e]})}(i);n("f25e");var o=n("2877"),s=Object(o["a"])(a["default"],r["a"],r["b"],!1,null,"747c530e",null);t["default"]=s.exports},"9c5e":function(e,t,n){"use strict";var r=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",[n("el-button",{on:{click:function(t){return e.ShowJsonDialog(e.objc)}}},[e._v("JSON")]),n("el-drawer",{ref:"drawer",attrs:{title:"JSON textarea!","before-close":e.handleClose,visible:e.dialog,direction:"ltr","custom-class":"demo-drawer"},on:{"update:visible":function(t){e.dialog=t}}},[n("el-input",{attrs:{type:"textarea",autosize:{minRows:10,maxRows:30},placeholder:"Please input json data!"},model:{value:e.textarea,callback:function(t){e.textarea=t},expression:"textarea"}})],1)],1)},a=[];n.d(t,"a",function(){return r}),n.d(t,"b",function(){return a})},"9d7e0":function(e,t,n){"use strict";var r=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",[e.checkResult?n("div",[n("el-col",{attrs:{span:2}},[n("el-tooltip",{attrs:{effect:"dark",content:"Please perfect liveId in this typeId!",placement:"bottom-start"}},[n("div",{staticClass:"grid-content bg-Danger",staticStyle:{"text-align":"center","line-height":"46px"}},[e._v(e._s(e.metaName))])])],1)],1):n("div",[n("el-col",{attrs:{span:2}},[n("div",{staticClass:"grid-content bg-purple",staticStyle:{"text-align":"center","line-height":"46px"}},[e._v(e._s(e.metaName))])])],1)])},a=[];n.d(t,"a",function(){return r}),n.d(t,"b",function(){return a})},ac4d:function(e,t,n){n("3a72")("asyncIterator")},bb3c:function(e,t,n){"use strict";var r=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",[n("el-row",[n("el-col",{attrs:{span:2}},[e._v("relyLives")]),n("el-col",{attrs:{span:2}},[n("el-button",{attrs:{disabled:e.disable},on:{click:function(t){return e.addObeject()}}},[e._v("add")])],1),n("el-col",{attrs:{span:3}},[n("json-rely-button",{attrs:{objc:e.parsedData},model:{value:e.parsedData,callback:function(t){e.parsedData=t},expression:"parsedData"}})],1)],1),e._l(e.parsedData,function(t,r){return n("el-row",{model:{value:e.parsedData,callback:function(t){e.parsedData=t},expression:"parsedData"}},[n("el-col",{attrs:{span:6,offset:2}},[n("el-input",{attrs:{type:"text"},model:{value:t.name,callback:function(n){e.$set(t,"name",n)},expression:"ob.name"}})],1),n("el-col",{attrs:{span:2}},[n("el-dropdown",{attrs:{trigger:"click"}},[n("span",{staticClass:"el-dropdown-link",staticStyle:{"text-align":"center","line-height":"46px"}},[e._v("\n        Select"),n("i",{staticClass:"el-icon-arrow-down el-icon--right"})]),n("el-dropdown-menu",{attrs:{slot:"dropdown"},slot:"dropdown"},e._l(e.$root.Configs,function(r){return n("div",e._l(r.lives,function(r){return n("div",[n("el-dropdown-item",{nativeOn:{click:function(n){return e.changeItem(t,r)}}},[e._v(e._s(r.name)+":"+e._s(r.liveId))])],1)}),0)}),0)],1)],1),n("el-col",{attrs:{span:10}},[n("el-input",{attrs:{type:"text"},model:{value:t.remark,callback:function(n){e.$set(t,"remark",n)},expression:"ob.remark"}})],1),n("el-col",{attrs:{span:3}},[n("el-button",{on:{click:function(t){return e.removeObject(r)}}},[e._v("remove")])],1)],1)})],2)},a=[];n.d(t,"a",function(){return r}),n.d(t,"b",function(){return a})},c034:function(e,t,n){"use strict";var r=function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("div",[n("el-collapse",{model:{value:e.activeTypes,callback:function(t){e.activeTypes=t},expression:"activeTypes"}},e._l(e.$root.Configs,function(t,r){return n("el-row",["not-exist"==t.metaData.flag?n("div",[n("el-col",{attrs:{span:1}},[n("el-tooltip",{attrs:{effect:"dark",content:"This typeId is not exist in dots!",placement:"bottom-start"}},[n("div",{staticClass:"grid-content bg-warning",staticStyle:{"text-align":"center","line-height":"46px"}},[e._v(e._s(r+1))])])],1)],1):n("div",[n("el-col",{attrs:{span:1}},[n("div",{staticClass:"grid-content bg-purple",staticStyle:{"text-align":"center","line-height":"46px"}},[e._v(e._s(r+1))])])],1),n("check-live-id",{attrs:{metaName:t.metaData.name,lives:t.lives}}),n("el-col",{attrs:{span:18}},[n("el-collapse-item",{attrs:{title:t.metaData.typeId,name:r}},e._l(t.lives,function(a,i){return n("el-row",[n("el-col",{attrs:{span:2}},[n("div",{staticClass:"grid-content bg-purple",staticStyle:{"text-align":"center","line-height":"46px"}},[e._v(e._s(a.name))])]),n("el-col",{attrs:{span:17}},[n("el-collapse-item",{attrs:{title:a.liveId,name:r+" "+i}},[n("el-row",[n("el-col",{attrs:{span:2}},[n("label",[e._v("name")])]),n("el-col",{attrs:{span:15}},[n("el-input",{attrs:{type:"text",placeholder:"Name"},model:{value:a.name,callback:function(t){e.$set(a,"name",t)},expression:"live.name"}})],1)],1),n("el-row",[n("el-col",{attrs:{span:2}},[n("label",[e._v("liveId")])]),n("el-col",{attrs:{span:15}},[n("el-input",{attrs:{type:"text",placeholder:"LiveId"},model:{value:a.liveId,callback:function(t){e.$set(a,"liveId",t)},expression:"live.liveId"}})],1),n("el-col",{attrs:{span:4}},[n("el-button",{on:{click:function(t){return e.UuidGenerator(a)}}},[e._v("Generate Live Id")])],1)],1),n("rely-lives-editor",{attrs:{objData:a},model:{value:a.relyLives,callback:function(t){e.$set(a,"relyLives",t)},expression:"live.relyLives"}}),a.json?n("el-row",[n("el-col",{attrs:{span:20}},[n("el-collapse-item",{attrs:{title:"Extend Config for live",name:r+","+i}},[n("extend-config-editor",{attrs:{objData:a.json},model:{value:a.json,callback:function(t){e.$set(a,"json",t)},expression:"live.json"}})],1)],1),n("el-col",{attrs:{span:4}},[n("json-button",{attrs:{objc:a.json},model:{value:a.json,callback:function(t){e.$set(a,"json",t)},expression:"live.json"}})],1)],1):e._e()],1)],1),n("el-col",{attrs:{span:2}},[n("json-button",{attrs:{objc:t.lives[i]},model:{value:t.lives[i],callback:function(n){e.$set(t.lives,i,n)},expression:"config.lives[index2]"}})],1),n("el-col",{attrs:{span:3}},[n("el-button",{attrs:{disabled:t.lives.length<=1},on:{click:function(n){return e.RemoveObject(t,i)}}},[e._v("Remove Live")])],1)],1)}),1)],1),n("el-col",{attrs:{span:3}},[n("el-row",[n("el-button",{attrs:{size:"mini"},on:{click:function(n){return e.showDialog(t.metaData.typeId,t.lives)}}},[e._v("Load By Config")])],1),n("el-row",[n("span",[n("el-button",{attrs:{size:"mini"},on:{click:function(n){return e.AddObject(t,t.metaData.typeId)}}},[e._v("Add Live")])],1),n("span",[n("el-button",{attrs:{size:"mini"},on:{click:function(t){return e.removeType(r)}}},[e._v("remove type")])],1)])],1)],1)}),1),n("el-dialog",{attrs:{title:"load by config",visible:e.dialogVisible,width:"40%"},on:{"update:visible":function(t){e.dialogVisible=t}}},[n("span",{staticStyle:{"text-align":"center"}},[n("el-upload",{staticClass:"upload-demo",attrs:{action:"","http-request":e.uploadSectionFile,drag:"",limit:1}},[n("i",{staticClass:"el-icon-upload"}),n("div",{staticClass:"el-upload__text"},[e._v("Drag and drop files here, or "),n("em",[e._v("click to upload")])]),n("div",{staticClass:"el-upload__tip",attrs:{slot:"tip"},slot:"tip"},[e._v("Can only upload json files")])])],1),n("span",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[n("el-button",{on:{click:function(t){e.dialogVisible=!1}}},[e._v("Cancel")]),n("el-button",{attrs:{type:"primary"},on:{click:function(t){return e.handleConfirm()}}},[e._v("Ok")])],1)])],1)},a=[];n.d(t,"a",function(){return r}),n.d(t,"b",function(){return a})},c39b:function(e,t,n){"use strict";var r=n("288e");n("1c01"),Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var a=r(n("2b0e")),i=a.default.extend({name:"checkLiveId",props:{metaName:{type:String,required:!0},lives:{type:Array,required:!0}},watch:{lives:{handler:function(e,t){for(var n=!1,r=0,a=e.length;r<a;r++)if(""===e[r].liveId){this.checkResult=!0,n=!0;break}n||(this.checkResult=!1)},immediate:!0,deep:!0}},data:function(){return{checkResult:!1}}});t.default=i},c437:function(e,t,n){var r,a,i=n("e1f4"),o=n("2366"),s=0,l=0;function c(e,t,n){var c=t&&n||0,u=t||[];e=e||{};var f=e.node||r,d=void 0!==e.clockseq?e.clockseq:a;if(null==f||null==d){var v=i();null==f&&(f=r=[1|v[0],v[1],v[2],v[3],v[4],v[5]]),null==d&&(d=a=16383&(v[6]<<8|v[7]))}var p=void 0!==e.msecs?e.msecs:(new Date).getTime(),h=void 0!==e.nsecs?e.nsecs:l+1,b=p-s+(h-l)/1e4;if(b<0&&void 0===e.clockseq&&(d=d+1&16383),(b<0||p>s)&&void 0===e.nsecs&&(h=0),h>=1e4)throw new Error("uuid.v1(): Can't create more than 10M uuids/sec");s=p,l=h,a=d,p+=122192928e5;var m=(1e4*(268435455&p)+h)%4294967296;u[c++]=m>>>24&255,u[c++]=m>>>16&255,u[c++]=m>>>8&255,u[c++]=255&m;var y=p/4294967296*1e4&268435455;u[c++]=y>>>8&255,u[c++]=255&y,u[c++]=y>>>24&15|16,u[c++]=y>>>16&255,u[c++]=d>>>8|128,u[c++]=255&d;for(var g=0;g<6;++g)u[c+g]=f[g];return t||o(u)}e.exports=c},c9d2:function(e,t,n){"use strict";n.r(t);var r=n("fdb9"),a=n.n(r);for(var i in r)"default"!==i&&function(e){n.d(t,e,function(){return r[e]})}(i);t["default"]=a.a},ccc3:function(e,t,n){"use strict";var r=n("288e");n("1c01"),Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0,n("7f7f");var a=r(n("2b0e")),i=r(n("8967")),o=n("5146"),s=a.default.extend({name:"RelyLivesEditor",props:{objData:{type:Object,required:!0}},data:function(){return{parsedData:[],disable:!1}},components:{"json-rely-button":i.default},watch:{objData:{handler:function(e,t){this.parsedData=(0,o.jsonParseRely)(this.objData.relyLives)},immediate:!0},parsedData:{handler:function(e,t){this.checkKey(),this.$emit("input",(0,o.makeJsonRely)(e))},deep:!0}},methods:{checkKey:function(){for(var e=0;e<this.parsedData.length;++e)if("default"===this.parsedData[e].name)return void(this.disable=!0);this.disable=!1},addObeject:function(){var e={name:"default",remark:"please change default"};this.parsedData.push(e)},removeObject:function(e){this.parsedData.splice(e,1)},changeItem:function(e,t){t.name&&(e.name=t.name),t.liveId&&(e.remark=t.liveId)}}});t.default=s},d4c0:function(e,t,n){var r=n("0d58"),a=n("2621"),i=n("52a7");e.exports=function(e){var t=r(e),n=a.f;if(n){var o,s=n(e),l=i.f,c=0;while(s.length>c)l.call(e,o=s[c++])&&t.push(o)}return t}},d4c1:function(e,t,n){"use strict";var r=n("fd74"),a=n.n(r);a.a},e1f4:function(e,t){var n="undefined"!=typeof crypto&&crypto.getRandomValues&&crypto.getRandomValues.bind(crypto)||"undefined"!=typeof msCrypto&&"function"==typeof window.msCrypto.getRandomValues&&msCrypto.getRandomValues.bind(msCrypto);if(n){var r=new Uint8Array(16);e.exports=function(){return n(r),r}}else{var a=new Array(16);e.exports=function(){for(var e,t=0;t<16;t++)0===(3&t)&&(e=4294967296*Math.random()),a[t]=e>>>((3&t)<<3)&255;return a}}},f25e:function(e,t,n){"use strict";var r=n("0b6e"),a=n.n(r);a.a},fd74:function(e,t,n){},fdb9:function(e,t,n){"use strict";var r=n("288e");n("1c01"),Object.defineProperty(t,"__esModule",{value:!0}),t.default=void 0;var a=r(n("7618")),i=r(n("2b0e")),o=n("5146"),s=i.default.extend({name:"JsonRelyButton",props:{objc:{type:Array,required:!0}},data:function(){return{dialog:!1,schemaObject:{},textarea:""}},methods:{ShowJsonDialog:function(e){this.dialog=!0,this.objc=e;var t=(0,o.makeJsonRely)(e);this.textarea=JSON.stringify(t,null,4)},handleClose:function(e){try{if(this.textarea){var t=JSON.parse(this.textarea);this.inputCheck(t)?this.$emit("input",(0,o.jsonParseRely)(t)):this.$message.error("json text input error!")}else this.$message.error("json text input error!")}catch(n){this.$message.error("json text input error!")}finally{e()}},inputCheck:function(e){for(var t in e)if((0,a.default)(e[t])!==(0,a.default)(""))return!1;return!0}}});t.default=s}}]);