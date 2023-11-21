"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[8671],{3905:(e,t,n)=>{n.d(t,{Zo:()=>p,kt:()=>g});var r=n(67294);function i(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function l(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){i(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function a(e,t){if(null==e)return{};var n,r,i=function(e,t){if(null==e)return{};var n,r,i={},o=Object.keys(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||(i[n]=e[n]);return i}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(i[n]=e[n])}return i}var c=r.createContext({}),s=function(e){var t=r.useContext(c),n=t;return e&&(n="function"==typeof e?e(t):l(l({},t),e)),n},p=function(e){var t=s(e.components);return r.createElement(c.Provider,{value:t},e.children)},u="mdxType",d={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},h=r.forwardRef((function(e,t){var n=e.components,i=e.mdxType,o=e.originalType,c=e.parentName,p=a(e,["components","mdxType","originalType","parentName"]),u=s(n),h=i,g=u["".concat(c,".").concat(h)]||u[h]||d[h]||o;return n?r.createElement(g,l(l({ref:t},p),{},{components:n})):r.createElement(g,l({ref:t},p))}));function g(e,t){var n=arguments,i=t&&t.mdxType;if("string"==typeof e||i){var o=n.length,l=new Array(o);l[0]=h;var a={};for(var c in t)hasOwnProperty.call(t,c)&&(a[c]=t[c]);a.originalType=e,a[u]="string"==typeof e?e:i,l[1]=a;for(var s=2;s<o;s++)l[s]=n[s];return r.createElement.apply(null,l)}return r.createElement.apply(null,n)}h.displayName="MDXCreateElement"},85596:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>c,contentTitle:()=>l,default:()=>d,frontMatter:()=>o,metadata:()=>a,toc:()=>s});var r=n(87462),i=(n(67294),n(3905));const o={title:"Integration",sidebar_label:"Integration",sidebar_position:2,slug:"/ibc/light-clients/localhost/integration"},l="Integration",a={unversionedId:"light-clients/localhost/integration",id:"light-clients/localhost/integration",title:"Integration",description:"The 09-localhost light client module registers codec types within the core IBC module. This differs from other light client module implementations which are expected to register codec types using the AppModuleBasic interface.",source:"@site/docs/03-light-clients/02-localhost/02-integration.md",sourceDirName:"03-light-clients/02-localhost",slug:"/ibc/light-clients/localhost/integration",permalink:"/main/ibc/light-clients/localhost/integration",draft:!1,tags:[],version:"current",sidebarPosition:2,frontMatter:{title:"Integration",sidebar_label:"Integration",sidebar_position:2,slug:"/ibc/light-clients/localhost/integration"},sidebar:"defaultSidebar",previous:{title:"Overview",permalink:"/main/ibc/light-clients/localhost/overview"},next:{title:"ClientState",permalink:"/main/ibc/light-clients/localhost/client-state"}},c={},s=[],p={toc:s},u="wrapper";function d(e){let{components:t,...n}=e;return(0,i.kt)(u,(0,r.Z)({},p,n,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"integration"},"Integration"),(0,i.kt)("p",null,"The 09-localhost light client module registers codec types within the core IBC module. This differs from other light client module implementations which are expected to register codec types using the ",(0,i.kt)("inlineCode",{parentName:"p"},"AppModuleBasic")," interface."),(0,i.kt)("p",null,"The localhost client is added to the 02-client submodule param ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc-go/blob/v7.0.0/proto/ibc/core/client/v1/client.proto#L102"},(0,i.kt)("inlineCode",{parentName:"a"},"allowed_clients"))," by default in ibc-go."),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},"var (\n  // DefaultAllowedClients are the default clients for the AllowedClients parameter.\n  DefaultAllowedClients = []string{exported.Solomachine, exported.Tendermint, exported.Localhost}\n)\n")))}d.isMDXComponent=!0}}]);