"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[9860],{3905:(e,t,i)=>{i.d(t,{Zo:()=>m,kt:()=>h});var n=i(67294);function r(e,t,i){return t in e?Object.defineProperty(e,t,{value:i,enumerable:!0,configurable:!0,writable:!0}):e[t]=i,e}function a(e,t){var i=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),i.push.apply(i,n)}return i}function o(e){for(var t=1;t<arguments.length;t++){var i=null!=arguments[t]?arguments[t]:{};t%2?a(Object(i),!0).forEach((function(t){r(e,t,i[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(i)):a(Object(i)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(i,t))}))}return e}function s(e,t){if(null==e)return{};var i,n,r=function(e,t){if(null==e)return{};var i,n,r={},a=Object.keys(e);for(n=0;n<a.length;n++)i=a[n],t.indexOf(i)>=0||(r[i]=e[i]);return r}(e,t);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(n=0;n<a.length;n++)i=a[n],t.indexOf(i)>=0||Object.prototype.propertyIsEnumerable.call(e,i)&&(r[i]=e[i])}return r}var l=n.createContext({}),c=function(e){var t=n.useContext(l),i=t;return e&&(i="function"==typeof e?e(t):o(o({},t),e)),i},m=function(e){var t=c(e.components);return n.createElement(l.Provider,{value:t},e.children)},p="mdxType",u={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},d=n.forwardRef((function(e,t){var i=e.components,r=e.mdxType,a=e.originalType,l=e.parentName,m=s(e,["components","mdxType","originalType","parentName"]),p=c(i),d=r,h=p["".concat(l,".").concat(d)]||p[d]||u[d]||a;return i?n.createElement(h,o(o({ref:t},m),{},{components:i})):n.createElement(h,o({ref:t},m))}));function h(e,t){var i=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var a=i.length,o=new Array(a);o[0]=d;var s={};for(var l in t)hasOwnProperty.call(t,l)&&(s[l]=t[l]);s.originalType=e,s[p]="string"==typeof e?e:r,o[1]=s;for(var c=2;c<a;c++)o[c]=i[c];return n.createElement.apply(null,o)}return n.createElement.apply(null,i)}d.displayName="MDXCreateElement"},61664:(e,t,i)=>{i.r(t),i.d(t,{assets:()=>l,contentTitle:()=>o,default:()=>u,frontMatter:()=>a,metadata:()=>s,toc:()=>c});var n=i(87462),r=(i(67294),i(3905));const a={title:"Overview",sidebar_label:"Overview",sidebar_position:1,slug:"/ibc/light-clients/wasm/overview"},o="08-wasm",s={unversionedId:"light-clients/wasm/overview",id:"light-clients/wasm/overview",title:"Overview",description:"Overview",source:"@site/docs/03-light-clients/04-wasm/01-overview.md",sourceDirName:"03-light-clients/04-wasm",slug:"/ibc/light-clients/wasm/overview",permalink:"/main/ibc/light-clients/wasm/overview",draft:!1,tags:[],version:"current",sidebarPosition:1,frontMatter:{title:"Overview",sidebar_label:"Overview",sidebar_position:1,slug:"/ibc/light-clients/wasm/overview"},sidebar:"defaultSidebar",previous:{title:"State Transitions",permalink:"/main/ibc/light-clients/solomachine/state_transitions"},next:{title:"Concepts",permalink:"/main/ibc/light-clients/wasm/concepts"}},l={},c=[{value:"Overview",id:"overview",level:2},{value:"Context",id:"context",level:3},{value:"Motivation",id:"motivation",level:3},{value:"Use cases",id:"use-cases",level:3}],m={toc:c},p="wrapper";function u(e){let{components:t,...i}=e;return(0,r.kt)(p,(0,n.Z)({},m,i,{components:t,mdxType:"MDXLayout"}),(0,r.kt)("h1",{id:"08-wasm"},(0,r.kt)("inlineCode",{parentName:"h1"},"08-wasm")),(0,r.kt)("h2",{id:"overview"},"Overview"),(0,r.kt)("p",null,"Learn about the ",(0,r.kt)("inlineCode",{parentName:"p"},"08-wasm")," light client proxy module. {synopsis}"),(0,r.kt)("h3",{id:"context"},"Context"),(0,r.kt)("p",null,"Traditionally, light clients used by ibc-go have been implemented only in Go, and since ibc-go v7 (with the release of the 02-client refactor), they are ",(0,r.kt)("a",{parentName:"p",href:"../../../architecture/adr-010-light-clients-as-sdk-modules.md"},"first-class Cosmos SDK modules"),". This means that updating existing light client implementations or adding support for new light clients is a multi-step, time-consuming process involving on-chain governance: it is necessary to modify the codebase of ibc-go (if the light client is part of its codebase), re-build chains' binaries, pass a governance proposal and have validators upgrade their nodes. "),(0,r.kt)("h3",{id:"motivation"},"Motivation"),(0,r.kt)("p",null,"To break the limitation of being able to write light client implementations only in Go, the ",(0,r.kt)("inlineCode",{parentName:"p"},"08-wasm")," adds support to run light clients written in a Wasm-compilable language. The light client byte code implements the entry points of a ",(0,r.kt)("a",{parentName:"p",href:"https://docs.cosmwasm.com/docs/"},"CosmWasm")," smart contract, and runs inside a Wasm VM. The ",(0,r.kt)("inlineCode",{parentName:"p"},"08-wasm")," module exposes a proxy light client interface that routes incoming messages to the appropriate handler function, inside the Wasm VM, for execution."),(0,r.kt)("p",null,"Adding a new light client to a chain is just as simple as submitting a governance proposal with the message that stores the byte code of the light client contract. No coordinated upgrade is needed. When the governance proposal passes and the message is executed, the contract is ready to be instantiated upon receiving a relayer-submitted ",(0,r.kt)("inlineCode",{parentName:"p"},"MsgCreateClient"),". The process of creating a Wasm light client is the same as with a regular light client implemented in Go."),(0,r.kt)("h3",{id:"use-cases"},"Use cases"),(0,r.kt)("ul",null,(0,r.kt)("li",{parentName:"ul"},"Development of light clients for non-Cosmos ecosystem chains: state machines in other ecosystems are, in many cases, implemented in Rust, and thus there are probably libraries used in their light client implementations for which there is no equivalent in Go. This makes the development of a light client in Go very difficult, but relatively simple to do it in Rust. Therefore, writing a CosmWasm smart contract in Rust that implements the light client algorithm becomes a lower effort.")))}u.isMDXComponent=!0}}]);