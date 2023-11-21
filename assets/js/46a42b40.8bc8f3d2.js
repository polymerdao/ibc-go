"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[8240],{3905:(e,n,t)=>{t.d(n,{Zo:()=>p,kt:()=>h});var a=t(67294);function r(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function o(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);n&&(a=a.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,a)}return t}function i(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?o(Object(t),!0).forEach((function(n){r(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):o(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function s(e,n){if(null==e)return{};var t,a,r=function(e,n){if(null==e)return{};var t,a,r={},o=Object.keys(e);for(a=0;a<o.length;a++)t=o[a],n.indexOf(t)>=0||(r[t]=e[t]);return r}(e,n);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(a=0;a<o.length;a++)t=o[a],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(r[t]=e[t])}return r}var l=a.createContext({}),c=function(e){var n=a.useContext(l),t=n;return e&&(t="function"==typeof e?e(n):i(i({},n),e)),t},p=function(e){var n=c(e.components);return a.createElement(l.Provider,{value:n},e.children)},u="mdxType",d={inlineCode:"code",wrapper:function(e){var n=e.children;return a.createElement(a.Fragment,{},n)}},m=a.forwardRef((function(e,n){var t=e.components,r=e.mdxType,o=e.originalType,l=e.parentName,p=s(e,["components","mdxType","originalType","parentName"]),u=c(t),m=r,h=u["".concat(l,".").concat(m)]||u[m]||d[m]||o;return t?a.createElement(h,i(i({ref:n},p),{},{components:t})):a.createElement(h,i({ref:n},p))}));function h(e,n){var t=arguments,r=n&&n.mdxType;if("string"==typeof e||r){var o=t.length,i=new Array(o);i[0]=m;var s={};for(var l in n)hasOwnProperty.call(n,l)&&(s[l]=n[l]);s.originalType=e,s[u]="string"==typeof e?e:r,i[1]=s;for(var c=2;c<o;c++)i[c]=t[c];return a.createElement.apply(null,i)}return a.createElement.apply(null,t)}m.displayName="MDXCreateElement"},3248:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>l,contentTitle:()=>i,default:()=>d,frontMatter:()=>o,metadata:()=>s,toc:()=>c});var a=t(87462),r=(t(67294),t(3905));const o={title:"Authentication Modules",sidebar_label:"Authentication Modules",sidebar_position:1,slug:"/apps/interchain-accounts/legacy/auth-modules"},i="Building an authentication module",s={unversionedId:"apps/interchain-accounts/legacy/auth-modules",id:"version-v6.2.x/apps/interchain-accounts/legacy/auth-modules",title:"Authentication Modules",description:"Deprecation Notice",source:"@site/versioned_docs/version-v6.2.x/02-apps/02-interchain-accounts/09-legacy/01-auth-modules.md",sourceDirName:"02-apps/02-interchain-accounts/09-legacy",slug:"/apps/interchain-accounts/legacy/auth-modules",permalink:"/v6/apps/interchain-accounts/legacy/auth-modules",draft:!1,tags:[],version:"v6.2.x",sidebarPosition:1,frontMatter:{title:"Authentication Modules",sidebar_label:"Authentication Modules",sidebar_position:1,slug:"/apps/interchain-accounts/legacy/auth-modules"},sidebar:"defaultSidebar",previous:{title:"Active Channels",permalink:"/v6/apps/interchain-accounts/active-channels"},next:{title:"Integration",permalink:"/v6/apps/interchain-accounts/legacy/integration"}},l={},c=[{value:"Deprecation Notice",id:"deprecation-notice",level:2},{value:"<code>IBCModule</code> implementation",id:"ibcmodule-implementation",level:2},{value:"<code>OnAcknowledgementPacket</code>",id:"onacknowledgementpacket",level:2},{value:"Integration into <code>app.go</code> file",id:"integration-into-appgo-file",level:2}],p={toc:c},u="wrapper";function d(e){let{components:n,...t}=e;return(0,r.kt)(u,(0,a.Z)({},p,t,{components:n,mdxType:"MDXLayout"}),(0,r.kt)("h1",{id:"building-an-authentication-module"},"Building an authentication module"),(0,r.kt)("h2",{id:"deprecation-notice"},"Deprecation Notice"),(0,r.kt)("p",null,(0,r.kt)("strong",{parentName:"p"},"This document is deprecated and will be removed in future releases"),"."),(0,r.kt)("admonition",{title:"Synopsis",type:"note"},(0,r.kt)("p",{parentName:"admonition"},"Authentication modules play the role of the ",(0,r.kt)("inlineCode",{parentName:"p"},"Base Application")," as described in ",(0,r.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc/tree/master/spec/app/ics-030-middleware"},"ICS-30 IBC Middleware"),", and enable application developers to perform custom logic when working with the Interchain Accounts controller API. ")),(0,r.kt)("p",null,"The controller submodule is used for account registration and packet sending. It executes only logic required of all controllers of interchain accounts. The type of authentication used to manage the interchain accounts remains unspecified. There may exist many different types of authentication which are desirable for different use cases. Thus the purpose of the authentication module is to wrap the controller submodule with custom authentication logic."),(0,r.kt)("p",null,"In ibc-go, authentication modules are connected to the controller chain via a middleware stack. The controller submodule is implemented as ",(0,r.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc/tree/master/spec/app/ics-030-middleware"},"middleware")," and the authentication module is connected to the controller submodule as the base application of the middleware stack. To implement an authentication module, the ",(0,r.kt)("inlineCode",{parentName:"p"},"IBCModule")," interface must be fulfilled. By implementing the controller submodule as middleware, any amount of authentication modules can be created and connected to the controller submodule without writing redundant code. "),(0,r.kt)("p",null,"The authentication module must:"),(0,r.kt)("ul",null,(0,r.kt)("li",{parentName:"ul"},"Authenticate interchain account owners."),(0,r.kt)("li",{parentName:"ul"},"Track the associated interchain account address for an owner."),(0,r.kt)("li",{parentName:"ul"},"Send packets on behalf of an owner (after authentication).")),(0,r.kt)("blockquote",null,(0,r.kt)("p",{parentName:"blockquote"},"Please note that since ibc-go v6 the channel capability is claimed by the controller submodule and therefore it is not required for authentication modules to claim the capability in the ",(0,r.kt)("inlineCode",{parentName:"p"},"OnChanOpenInit")," callback. When the authentication module sends packets on the channel created for the associated interchain account it can pass a ",(0,r.kt)("inlineCode",{parentName:"p"},"nil")," capability to the legacy function ",(0,r.kt)("inlineCode",{parentName:"p"},"SendTx")," of the controller keeper (see ",(0,r.kt)("a",{parentName:"p",href:"/v6/apps/interchain-accounts/legacy/keeper-api#sendtx"},"section ",(0,r.kt)("inlineCode",{parentName:"a"},"SendTx"))," below for mode information). ")),(0,r.kt)("h2",{id:"ibcmodule-implementation"},(0,r.kt)("inlineCode",{parentName:"h2"},"IBCModule")," implementation"),(0,r.kt)("p",null,"The following ",(0,r.kt)("inlineCode",{parentName:"p"},"IBCModule")," callbacks must be implemented with appropriate custom logic:"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-go"},"// OnChanOpenInit implements the IBCModule interface\nfunc (im IBCModule) OnChanOpenInit(\n    ctx sdk.Context,\n    order channeltypes.Order,\n    connectionHops []string,\n    portID string,\n    channelID string,\n    chanCap *capabilitytypes.Capability,\n    counterparty channeltypes.Counterparty,\n    version string,\n) (string, error) {\n    // since ibc-go v6 the authentication module *must not* claim the channel capability on OnChanOpenInit\n\n    // perform custom logic\n\n    return version, nil\n}\n\n// OnChanOpenAck implements the IBCModule interface\nfunc (im IBCModule) OnChanOpenAck(\n    ctx sdk.Context,\n    portID,\n    channelID string,\n    counterpartyVersion string,\n) error {\n    // perform custom logic\n\n    return nil\n}\n\n// OnChanCloseConfirm implements the IBCModule interface\nfunc (im IBCModule) OnChanCloseConfirm(\n    ctx sdk.Context,\n    portID,\n    channelID string,\n) error {\n    // perform custom logic\n\n    return nil\n}\n\n// OnAcknowledgementPacket implements the IBCModule interface\nfunc (im IBCModule) OnAcknowledgementPacket(\n    ctx sdk.Context,\n    packet channeltypes.Packet,\n    acknowledgement []byte,\n    relayer sdk.AccAddress,\n) error {\n    // perform custom logic\n\n    return nil\n}\n\n// OnTimeoutPacket implements the IBCModule interface.\nfunc (im IBCModule) OnTimeoutPacket(\n    ctx sdk.Context,\n    packet channeltypes.Packet,\n    relayer sdk.AccAddress,\n) error {\n    // perform custom logic\n\n    return nil\n}\n")),(0,r.kt)("p",null,"The following functions must be defined to fulfill the ",(0,r.kt)("inlineCode",{parentName:"p"},"IBCModule")," interface, but they will never be called by the controller submodule so they may error or panic. That is because in Interchain Accounts, the channel handshake is always initiated on the controller chain and packets are always sent to the host chain and never to the controller chain."),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-go"},'// OnChanOpenTry implements the IBCModule interface\nfunc (im IBCModule) OnChanOpenTry(\n    ctx sdk.Context,\n    order channeltypes.Order,\n    connectionHops []string,\n    portID,\n    channelID string,\n    chanCap *capabilitytypes.Capability,\n    counterparty channeltypes.Counterparty,\n    counterpartyVersion string,\n) (string, error) {\n    panic("UNIMPLEMENTED")\n}\n\n// OnChanOpenConfirm implements the IBCModule interface\nfunc (im IBCModule) OnChanOpenConfirm(\n    ctx sdk.Context,\n    portID,\n    channelID string,\n) error {\n    panic("UNIMPLEMENTED")\n}\n\n// OnChanCloseInit implements the IBCModule interface\nfunc (im IBCModule) OnChanCloseInit(\n    ctx sdk.Context,\n    portID,\n    channelID string,\n) error {\n    panic("UNIMPLEMENTED")\n}\n\n// OnRecvPacket implements the IBCModule interface. A successful acknowledgement\n// is returned if the packet data is successfully decoded and the receive application\n// logic returns without error.\nfunc (im IBCModule) OnRecvPacket(\n    ctx sdk.Context,\n    packet channeltypes.Packet,\n    relayer sdk.AccAddress,\n) ibcexported.Acknowledgement {\n    panic("UNIMPLEMENTED")\n}\n')),(0,r.kt)("h2",{id:"onacknowledgementpacket"},(0,r.kt)("inlineCode",{parentName:"h2"},"OnAcknowledgementPacket")),(0,r.kt)("p",null,"Controller chains will be able to access the acknowledgement written into the host chain state once a relayer relays the acknowledgement.\nThe acknowledgement bytes contain either the response of the execution of the message(s) on the host chain or an error. They will be passed to the auth module via the ",(0,r.kt)("inlineCode",{parentName:"p"},"OnAcknowledgementPacket")," callback. Auth modules are expected to know how to decode the acknowledgement. "),(0,r.kt)("p",null,"If the controller chain is connected to a host chain using the host module on ibc-go, it may interpret the acknowledgement bytes as follows:"),(0,r.kt)("p",null,"Begin by unmarshaling the acknowledgement into ",(0,r.kt)("inlineCode",{parentName:"p"},"sdk.TxMsgData"),":"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-go"},"var ack channeltypes.Acknowledgement\nif err := channeltypes.SubModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {\n    return err\n}\n\ntxMsgData := &sdk.TxMsgData{}\nif err := proto.Unmarshal(ack.GetResult(), txMsgData); err != nil {\n    return err\n}\n")),(0,r.kt)("p",null,"If the ",(0,r.kt)("inlineCode",{parentName:"p"},"txMsgData.Data")," field is non nil, the host chain is using SDK version <= v0.45.\nThe auth module should interpret the ",(0,r.kt)("inlineCode",{parentName:"p"},"txMsgData.Data")," as follows:"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-go"},"switch len(txMsgData.Data) {\ncase 0:\n    // see documentation below for SDK 0.46.x or greater\ndefault:\n    for _, msgData := range txMsgData.Data {\n        if err := handler(msgData); err != nil {\n            return err\n        }\n    }\n...\n}            \n")),(0,r.kt)("p",null,"A handler will be needed to interpret what actions to perform based on the message type sent.\nA router could be used, or more simply a switch statement."),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-go"},"func handler(msgData sdk.MsgData) error {\nswitch msgData.MsgType {\ncase sdk.MsgTypeURL(&banktypes.MsgSend{}):\n    msgResponse := &banktypes.MsgSendResponse{}\n    if err := proto.Unmarshal(msgData.Data, msgResponse}; err != nil {\n        return err\n    }\n\n    handleBankSendMsg(msgResponse)\n\ncase sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}):\n    msgResponse := &stakingtypes.MsgDelegateResponse{}\n    if err := proto.Unmarshal(msgData.Data, msgResponse}; err != nil {\n        return err\n    }\n\n    handleStakingDelegateMsg(msgResponse)\n\ncase sdk.MsgTypeURL(&transfertypes.MsgTransfer{}):\n    msgResponse := &transfertypes.MsgTransferResponse{}\n    if err := proto.Unmarshal(msgData.Data, msgResponse}; err != nil {\n        return err\n    }\n\n    handleIBCTransferMsg(msgResponse)\n \ndefault:\n    return\n}\n")),(0,r.kt)("p",null,"If the ",(0,r.kt)("inlineCode",{parentName:"p"},"txMsgData.Data")," is empty, the host chain is using SDK version > v0.45.\nThe auth module should interpret the ",(0,r.kt)("inlineCode",{parentName:"p"},"txMsgData.Responses")," as follows:"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-go"},"...\n// switch statement from above\ncase 0:\n    for _, any := range txMsgData.MsgResponses {\n        if err := handleAny(any); err != nil {\n            return err\n        }\n    }\n}\n")),(0,r.kt)("p",null,"A handler will be needed to interpret what actions to perform based on the type URL of the Any.\nA router could be used, or more simply a switch statement.\nIt may be possible to deduplicate logic between ",(0,r.kt)("inlineCode",{parentName:"p"},"handler")," and ",(0,r.kt)("inlineCode",{parentName:"p"},"handleAny"),"."),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-go"},"func handleAny(any *codectypes.Any) error {\nswitch any.TypeURL {\ncase banktypes.MsgSend:\n    msgResponse, err := unpackBankMsgSendResponse(any)\n    if err != nil {\n        return err\n    }\n\n    handleBankSendMsg(msgResponse)\n\ncase stakingtypes.MsgDelegate:\n    msgResponse, err := unpackStakingDelegateResponse(any)\n    if err != nil {\n        return err\n    }\n\n    handleStakingDelegateMsg(msgResponse)\n\n    case transfertypes.MsgTransfer:\n    msgResponse, err := unpackIBCTransferMsgResponse(any)\n    if err != nil {\n        return err\n    }\n\n    handleIBCTransferMsg(msgResponse)\n \ndefault:\n    return\n}\n")),(0,r.kt)("h2",{id:"integration-into-appgo-file"},"Integration into ",(0,r.kt)("inlineCode",{parentName:"h2"},"app.go")," file"),(0,r.kt)("p",null,"To integrate the authentication module into your chain, please follow the steps outlined in ",(0,r.kt)("a",{parentName:"p",href:"/v6/apps/interchain-accounts/legacy/integration#example-integration"},(0,r.kt)("inlineCode",{parentName:"a"},"app.go")," integration"),"."))}d.isMDXComponent=!0}}]);