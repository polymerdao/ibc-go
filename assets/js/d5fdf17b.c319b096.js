"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[9834],{3905:(e,t,n)=>{n.d(t,{Zo:()=>m,kt:()=>d});var a=n(67294);function o(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function s(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?s(Object(n),!0).forEach((function(t){o(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):s(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function r(e,t){if(null==e)return{};var n,a,o=function(e,t){if(null==e)return{};var n,a,o={},s=Object.keys(e);for(a=0;a<s.length;a++)n=s[a],t.indexOf(n)>=0||(o[n]=e[n]);return o}(e,t);if(Object.getOwnPropertySymbols){var s=Object.getOwnPropertySymbols(e);for(a=0;a<s.length;a++)n=s[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(o[n]=e[n])}return o}var c=a.createContext({}),l=function(e){var t=a.useContext(c),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},m=function(e){var t=l(e.components);return a.createElement(c.Provider,{value:t},e.children)},p="mdxType",h={inlineCode:"code",wrapper:function(e){var t=e.children;return a.createElement(a.Fragment,{},t)}},g=a.forwardRef((function(e,t){var n=e.components,o=e.mdxType,s=e.originalType,c=e.parentName,m=r(e,["components","mdxType","originalType","parentName"]),p=l(n),g=o,d=p["".concat(c,".").concat(g)]||p[g]||h[g]||s;return n?a.createElement(d,i(i({ref:t},m),{},{components:n})):a.createElement(d,i({ref:t},m))}));function d(e,t){var n=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var s=n.length,i=new Array(s);i[0]=g;var r={};for(var c in t)hasOwnProperty.call(t,c)&&(r[c]=t[c]);r.originalType=e,r[p]="string"==typeof e?e:o,i[1]=r;for(var l=2;l<s;l++)i[l]=n[l];return a.createElement.apply(null,i)}return a.createElement.apply(null,n)}g.displayName="MDXCreateElement"},82486:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>c,contentTitle:()=>i,default:()=>h,frontMatter:()=>s,metadata:()=>r,toc:()=>l});var a=n(87462),o=(n(67294),n(3905));const s={title:"Governance",sidebar_label:"Governance",sidebar_position:5,slug:"/ibc/light-clients/wasm/governance"},i="Governance",r={unversionedId:"light-clients/wasm/governance",id:"light-clients/wasm/governance",title:"Governance",description:"Learn how to upload Wasm light client byte code on a chain, and how to migrate an existing Wasm light client contract.",source:"@site/docs/03-light-clients/04-wasm/05-governance.md",sourceDirName:"03-light-clients/04-wasm",slug:"/ibc/light-clients/wasm/governance",permalink:"/main/ibc/light-clients/wasm/governance",draft:!1,tags:[],version:"current",sidebarPosition:5,frontMatter:{title:"Governance",sidebar_label:"Governance",sidebar_position:5,slug:"/ibc/light-clients/wasm/governance"},sidebar:"defaultSidebar",previous:{title:"Messages",permalink:"/main/ibc/light-clients/wasm/messages"},next:{title:"Events",permalink:"/main/ibc/light-clients/wasm/events"}},c={},l=[{value:"Setting an authority",id:"setting-an-authority",level:2},{value:"Storing new Wasm light client byte code",id:"storing-new-wasm-light-client-byte-code",level:2},{value:"Migrating an existing Wasm light client contract",id:"migrating-an-existing-wasm-light-client-contract",level:2},{value:"Removing an existing checksum",id:"removing-an-existing-checksum",level:2}],m={toc:l},p="wrapper";function h(e){let{components:t,...n}=e;return(0,o.kt)(p,(0,a.Z)({},m,n,{components:t,mdxType:"MDXLayout"}),(0,o.kt)("h1",{id:"governance"},"Governance"),(0,o.kt)("p",null,"Learn how to upload Wasm light client byte code on a chain, and how to migrate an existing Wasm light client contract. {synopsis}"),(0,o.kt)("h2",{id:"setting-an-authority"},"Setting an authority"),(0,o.kt)("p",null,"Both the storage of Wasm light client byte code as well as the migration of an existing Wasm light client contract are permissioned (i.e. only allowed to an authority such as governance). The designated authority is specified when instantiating ",(0,o.kt)("inlineCode",{parentName:"p"},"08-wasm"),"'s keeper: both ",(0,o.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc-go/blob/c95c22f45cb217d27aca2665af9ac60b0d2f3a0c/modules/light-clients/08-wasm/keeper/keeper.go#L33-L38"},(0,o.kt)("inlineCode",{parentName:"a"},"NewKeeperWithVM"))," and ",(0,o.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc-go/blob/c95c22f45cb217d27aca2665af9ac60b0d2f3a0c/modules/light-clients/08-wasm/keeper/keeper.go#L52-L57"},(0,o.kt)("inlineCode",{parentName:"a"},"NewKeeperWithConfig"))," constructor functions accept an ",(0,o.kt)("inlineCode",{parentName:"p"},"authority")," argument that must be the address of the authorized actor. For example, in ",(0,o.kt)("inlineCode",{parentName:"p"},"app.go"),", when instantiating the keeper, you can pass the address of the governance module:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-go"},'// app.go\nimport (\n  ...\n  "github.com/cosmos/cosmos-sdk/runtime"\n  authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"\n  govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"\n\n  wasmkeeper "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/keeper"\n  wasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/types"\n  ...\n)\n\n// app.go\napp.WasmClientKeeper = wasmkeeper.NewKeeperWithVM(\n  appCodec,\n  runtime.NewKVStoreService(keys[wasmtypes.StoreKey]),\n  app.IBCKeeper.ClientKeeper,\n    authtypes.NewModuleAddress(govtypes.ModuleName).String(), // authority\n  wasmVM,\n)\n')),(0,o.kt)("h2",{id:"storing-new-wasm-light-client-byte-code"},"Storing new Wasm light client byte code"),(0,o.kt)("p",null," If governance is the allowed authority, the governance v1 proposal that needs to be submitted to upload a new light client contract should contain the message ",(0,o.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc-go/blob/f822b4fa7932a657420aba219c563e06c4465221/proto/ibc/lightclients/wasm/v1/tx.proto#L16-L23"},(0,o.kt)("inlineCode",{parentName:"a"},"MsgStoreCode"))," with the base64-encoded byte code of the Wasm contract. Use the following CLI command and JSON as an example:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-shell"},"simd tx gov submit-proposal <path/to/proposal.json> --from <key_or_address>\n")),(0,o.kt)("p",null,"where ",(0,o.kt)("inlineCode",{parentName:"p"},"proposal.json")," contains:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-json"},'{\n  "title": "Upload IBC Wasm light client",\n  "summary": "Upload wasm client",\n  "messages": [\n    {\n      "@type": "/ibc.lightclients.wasm.v1.MsgStoreCode",\n      "signer": "cosmos1...", // the authority address (e.g. the gov module account address)\n      "wasm_byte_code": "YWJ...PUB+" // standard base64 encoding of the Wasm contract byte code\n    }\n  ],\n  "metadata": "AQ==",\n  "deposit": "100stake"\n}\n')),(0,o.kt)("p",null,"To learn more about the ",(0,o.kt)("inlineCode",{parentName:"p"},"submit-proposal")," CLI command, please check out ",(0,o.kt)("a",{parentName:"p",href:"https://docs.cosmos.network/main/modules/gov#submit-proposal"},"the relevant section in Cosmos SDK documentation"),"."),(0,o.kt)("p",null,"Alternatively, the process of submitting the proposal may be simpler if you use the CLI command ",(0,o.kt)("inlineCode",{parentName:"p"},"store-code"),". This CLI command accepts as argument the file of the Wasm light client contract and takes care of constructing the proposal message with ",(0,o.kt)("inlineCode",{parentName:"p"},"MsgStoreCode")," and broadcasting it. See section ",(0,o.kt)("a",{parentName:"p",href:"/main/ibc/light-clients/wasm/client#store-code"},(0,o.kt)("inlineCode",{parentName:"a"},"store-code"))," for more information."),(0,o.kt)("h2",{id:"migrating-an-existing-wasm-light-client-contract"},"Migrating an existing Wasm light client contract"),(0,o.kt)("p",null,"If governance is the allowed authority, the governance v1 proposal that needs to be submitted to migrate an existing new Wasm light client contract should contain the message ",(0,o.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc-go/blob/729cb090951b1e996427b2258cf72c49787b885a/proto/ibc/lightclients/wasm/v1/tx.proto#L51-L63"},(0,o.kt)("inlineCode",{parentName:"a"},"MsgMigrateContract"))," with the checksum of the Wasm byte code to migrate to. Use the following CLI command and JSON as an example:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-shell"},"simd tx gov submit-proposal <path/to/proposal.json> --from <key_or_address>\n")),(0,o.kt)("p",null,"where ",(0,o.kt)("inlineCode",{parentName:"p"},"proposal.json")," contains:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-json"},'{\n  "title": "Migrate IBC Wasm light client",\n  "summary": "Migrate wasm client",\n  "messages": [\n    {\n      "@type": "/ibc.lightclients.wasm.v1.MsgMigrateContract",\n      "signer": "cosmos1...", // the authority address (e.g. the gov module account address)\n      "client_id": "08-wasm-1", // client identifier of the Wasm light client contract that will be migrated\n      "checksum": "a8ad...4dc0", // SHA-256 hash of the Wasm byte code to migrate to, previously stored with MsgStoreCode\n      "msg": "{}" // JSON-encoded message to be passed to the contract on migration\n    }\n  ],\n  "metadata": "AQ==",\n  "deposit": "100stake"\n}\n')),(0,o.kt)("p",null,"To learn more about the ",(0,o.kt)("inlineCode",{parentName:"p"},"submit-proposal")," CLI command, please check out ",(0,o.kt)("a",{parentName:"p",href:"https://docs.cosmos.network/main/modules/gov#submit-proposal"},"the relevant section in Cosmos SDK documentation"),"."),(0,o.kt)("h2",{id:"removing-an-existing-checksum"},"Removing an existing checksum"),(0,o.kt)("p",null,"If governance is the allowed authority, the governance v1 proposal that needs to be submitted to remove a specific checksum from the list of allowed checksums should contain the message ",(0,o.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc-go/blob/729cb090951b1e996427b2258cf72c49787b885a/proto/ibc/lightclients/wasm/v1/tx.proto#L38-L46"},(0,o.kt)("inlineCode",{parentName:"a"},"MsgRemoveChecksum"))," with the checksum (of a corresponding Wasm byte code). Use the following CLI command and JSON as an example:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-shell"},"simd tx gov submit-proposal <path/to/proposal.json> --from <key_or_address>\n")),(0,o.kt)("p",null,"where ",(0,o.kt)("inlineCode",{parentName:"p"},"proposal.json")," contains:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-json"},'{\n  "title": "Remove checksum of Wasm light client byte code",\n  "summary": "Remove checksum",\n  "messages": [\n    {\n      "@type": "/ibc.lightclients.wasm.v1.MsgRemoveChecksum",\n      "signer": "cosmos1...", // the authority address (e.g. the gov module account address)\n      "checksum": "a8ad...4dc0", // SHA-256 hash of the Wasm byte code that should be removed from the list of allowed checksums\n    }\n  ],\n  "metadata": "AQ==",\n  "deposit": "100stake"\n}\n')),(0,o.kt)("p",null,"To learn more about the ",(0,o.kt)("inlineCode",{parentName:"p"},"submit-proposal")," CLI command, please check out ",(0,o.kt)("a",{parentName:"p",href:"https://docs.cosmos.network/main/modules/gov#submit-proposal"},"the relevant section in Cosmos SDK documentation"),"."))}h.isMDXComponent=!0}}]);