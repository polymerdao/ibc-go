"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[9641],{3905:(e,n,a)=>{a.d(n,{Zo:()=>s,kt:()=>u});var t=a(67294);function i(e,n,a){return n in e?Object.defineProperty(e,n,{value:a,enumerable:!0,configurable:!0,writable:!0}):e[n]=a,e}function o(e,n){var a=Object.keys(e);if(Object.getOwnPropertySymbols){var t=Object.getOwnPropertySymbols(e);n&&(t=t.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),a.push.apply(a,t)}return a}function r(e){for(var n=1;n<arguments.length;n++){var a=null!=arguments[n]?arguments[n]:{};n%2?o(Object(a),!0).forEach((function(n){i(e,n,a[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(a)):o(Object(a)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(a,n))}))}return e}function l(e,n){if(null==e)return{};var a,t,i=function(e,n){if(null==e)return{};var a,t,i={},o=Object.keys(e);for(t=0;t<o.length;t++)a=o[t],n.indexOf(a)>=0||(i[a]=e[a]);return i}(e,n);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(t=0;t<o.length;t++)a=o[t],n.indexOf(a)>=0||Object.prototype.propertyIsEnumerable.call(e,a)&&(i[a]=e[a])}return i}var c=t.createContext({}),p=function(e){var n=t.useContext(c),a=n;return e&&(a="function"==typeof e?e(n):r(r({},n),e)),a},s=function(e){var n=p(e.components);return t.createElement(c.Provider,{value:n},e.children)},d="mdxType",m={inlineCode:"code",wrapper:function(e){var n=e.children;return t.createElement(t.Fragment,{},n)}},h=t.forwardRef((function(e,n){var a=e.components,i=e.mdxType,o=e.originalType,c=e.parentName,s=l(e,["components","mdxType","originalType","parentName"]),d=p(a),h=i,u=d["".concat(c,".").concat(h)]||d[h]||m[h]||o;return a?t.createElement(u,r(r({ref:n},s),{},{components:a})):t.createElement(u,r({ref:n},s))}));function u(e,n){var a=arguments,i=n&&n.mdxType;if("string"==typeof e||i){var o=a.length,r=new Array(o);r[0]=h;var l={};for(var c in n)hasOwnProperty.call(n,c)&&(l[c]=n[c]);l.originalType=e,l[d]="string"==typeof e?e:i,r[1]=l;for(var p=2;p<o;p++)r[p]=a[p];return t.createElement.apply(null,r)}return t.createElement.apply(null,a)}h.displayName="MDXCreateElement"},78109:(e,n,a)=>{a.r(n),a.d(n,{assets:()=>c,contentTitle:()=>r,default:()=>m,frontMatter:()=>o,metadata:()=>l,toc:()=>p});var t=a(87462),i=(a(67294),a(3905));const o={title:"IBC middleware",sidebar_label:"IBC middleware",sidebar_position:1,slug:"/ibc/middleware/develop"},r="IBC middleware",l={unversionedId:"ibc/middleware/develop",id:"version-v4.5.x/ibc/middleware/develop",title:"IBC middleware",description:"Learn how to write your own custom middleware to wrap an IBC application, and understand how to hook different middleware to IBC base applications to form different IBC application stacks",source:"@site/versioned_docs/version-v4.5.x/01-ibc/04-middleware/01-develop.md",sourceDirName:"01-ibc/04-middleware",slug:"/ibc/middleware/develop",permalink:"/v4/ibc/middleware/develop",draft:!1,tags:[],version:"v4.5.x",sidebarPosition:1,frontMatter:{title:"IBC middleware",sidebar_label:"IBC middleware",sidebar_position:1,slug:"/ibc/middleware/develop"}},c={},p=[{value:"Definitions",id:"definitions",level:2},{value:"Create a custom IBC middleware",id:"create-a-custom-ibc-middleware",level:2},{value:"Interfaces",id:"interfaces",level:3},{value:"Implement <code>IBCModule</code> interface and callbacks",id:"implement-ibcmodule-interface-and-callbacks",level:3},{value:"Handshake callbacks",id:"handshake-callbacks",level:3},{value:"<code>OnChanOpenInit</code>",id:"onchanopeninit",level:4},{value:"<code>OnChanOpenTry</code>",id:"onchanopentry",level:4},{value:"<code>OnChanOpenAck</code>",id:"onchanopenack",level:4},{value:"<code>OnChanOpenConfirm</code>",id:"onchanopenconfirm",level:3},{value:"<code>OnChanCloseInit</code>",id:"onchancloseinit",level:4},{value:"<code>OnChanCloseConfirm</code>",id:"onchancloseconfirm",level:4},{value:"Packet callbacks",id:"packet-callbacks",level:3},{value:"<code>OnRecvPacket</code>",id:"onrecvpacket",level:4},{value:"<code>OnAcknowledgementPacket</code>",id:"onacknowledgementpacket",level:4},{value:"<code>OnTimeoutPacket</code>",id:"ontimeoutpacket",level:4},{value:"ICS-4 wrappers",id:"ics-4-wrappers",level:3},{value:"<code>SendPacket</code>",id:"sendpacket",level:4},{value:"<code>WriteAcknowledgement</code>",id:"writeacknowledgement",level:4},{value:"<code>GetAppVersion</code>",id:"getappversion",level:4}],s={toc:p},d="wrapper";function m(e){let{components:n,...a}=e;return(0,i.kt)(d,(0,t.Z)({},s,a,{components:n,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"ibc-middleware"},"IBC middleware"),(0,i.kt)("admonition",{title:"Synopsis",type:"note"},(0,i.kt)("p",{parentName:"admonition"},"Learn how to write your own custom middleware to wrap an IBC application, and understand how to hook different middleware to IBC base applications to form different IBC application stacks")),(0,i.kt)("p",null,"This document serves as a guide for middleware developers who want to write their own middleware and for chain developers who want to use IBC middleware on their chains."),(0,i.kt)("p",null,"IBC applications are designed to be self-contained modules that implement their own application-specific logic through a set of interfaces with the core IBC handlers. These core IBC handlers, in turn, are designed to enforce the correctness properties of IBC (transport, authentication, ordering) while delegating all application-specific handling to the IBC application modules. However, there are cases where some functionality may be desired by many applications, yet not appropriate to place in core IBC."),(0,i.kt)("p",null,"Middleware allows developers to define the extensions as separate modules that can wrap over the base application. This middleware can thus perform its own custom logic, and pass data into the application so that it may run its logic without being aware of the middleware's existence. This allows both the application and the middleware to implement its own isolated logic while still being able to run as part of a single packet flow."),(0,i.kt)("admonition",{type:"note"},(0,i.kt)("h2",{parentName:"admonition",id:"pre-requisite-readings"},"Pre-requisite readings"),(0,i.kt)("ul",{parentName:"admonition"},(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("a",{parentName:"li",href:"/v4/ibc/overview"},"IBC Overview")),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("a",{parentName:"li",href:"/v4/ibc/integration"},"IBC Integration")),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("a",{parentName:"li",href:"/v4/ibc/apps/apps"},"IBC Application Developer Guide")))),(0,i.kt)("h2",{id:"definitions"},"Definitions"),(0,i.kt)("p",null,(0,i.kt)("inlineCode",{parentName:"p"},"Middleware"),": A self-contained module that sits between core IBC and an underlying IBC application during packet execution. All messages between core IBC and underlying application must flow through middleware, which may perform its own custom logic."),(0,i.kt)("p",null,(0,i.kt)("inlineCode",{parentName:"p"},"Underlying Application"),": An underlying application is the application that is directly connected to the middleware in question. This underlying application may itself be middleware that is chained to a base application."),(0,i.kt)("p",null,(0,i.kt)("inlineCode",{parentName:"p"},"Base Application"),": A base application is an IBC application that does not contain any middleware. It may be nested by 0 or multiple middleware to form an application stack."),(0,i.kt)("p",null,(0,i.kt)("inlineCode",{parentName:"p"},"Application Stack (or stack)"),": A stack is the complete set of application logic (middleware(s) + base application) that gets connected to core IBC. A stack may be just a base application, or it may be a series of middlewares that nest a base application."),(0,i.kt)("h2",{id:"create-a-custom-ibc-middleware"},"Create a custom IBC middleware"),(0,i.kt)("p",null,"IBC middleware will wrap over an underlying IBC application and sits between core IBC and the application. It has complete control in modifying any message coming from IBC to the application, and any message coming from the application to core IBC. Thus, middleware must be completely trusted by chain developers who wish to integrate them, however this gives them complete flexibility in modifying the application(s) they wrap."),(0,i.kt)("h3",{id:"interfaces"},"Interfaces"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},"// Middleware implements the ICS26 Module interface\ntype Middleware interface {\n    porttypes.IBCModule // middleware has acccess to an underlying application which may be wrapped by more middleware\n    ics4Wrapper: ICS4Wrapper // middleware has access to ICS4Wrapper which may be core IBC Channel Handler or a higher-level middleware that wraps this middleware.\n}\n")),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-typescript"},"// This is implemented by ICS4 and all middleware that are wrapping base application.\n// The base application will call `sendPacket` or `writeAcknowledgement` of the middleware directly above them\n// which will call the next middleware until it reaches the core IBC handler.\ntype ICS4Wrapper interface {\n    SendPacket(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.Packet) error\n    WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.Packet, ack exported.Acknowledgement) error\n    GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool)\n}\n")),(0,i.kt)("h3",{id:"implement-ibcmodule-interface-and-callbacks"},"Implement ",(0,i.kt)("inlineCode",{parentName:"h3"},"IBCModule")," interface and callbacks"),(0,i.kt)("p",null,"The ",(0,i.kt)("inlineCode",{parentName:"p"},"IBCModule")," is a struct that implements the ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc-go/blob/main/modules/core/05-port/types/module.go#L11-L106"},"ICS-26 interface (",(0,i.kt)("inlineCode",{parentName:"a"},"porttypes.IBCModule"),")"),". It is recommended to separate these callbacks into a separate file ",(0,i.kt)("inlineCode",{parentName:"p"},"ibc_module.go"),". As will be mentioned in the ",(0,i.kt)("a",{parentName:"p",href:"/v4/ibc/middleware/integration"},"integration section"),", this struct should be different than the struct that implements ",(0,i.kt)("inlineCode",{parentName:"p"},"AppModule")," in case the middleware maintains its own internal state and processes separate SDK messages."),(0,i.kt)("p",null,"The middleware must have access to the underlying application, and be called before during all ICS-26 callbacks. It may execute custom logic during these callbacks, and then call the underlying application's callback. Middleware ",(0,i.kt)("strong",{parentName:"p"},"may")," choose not to call the underlying application's callback at all. Though these should generally be limited to error cases."),(0,i.kt)("p",null,"In the case where the IBC middleware expects to speak to a compatible IBC middleware on the counterparty chain, they must use the channel handshake to negotiate the middleware version without interfering in the version negotiation of the underlying application."),(0,i.kt)("p",null,"Middleware accomplishes this by formatting the version in a JSON-encoded string containing the middleware version and the application version. The application version may as well be a JSON-encoded string, possibly including further middleware and app versions, if the application stack consists of multiple milddlewares wrapping a base application. The format of the version is specified in ICS-30 as the following:"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-json"},'{\n  "<middleware_version_key>": "<middleware_version_value>",\n  "app_version": "<application_version_value>"\n}\n')),(0,i.kt)("p",null,"The ",(0,i.kt)("inlineCode",{parentName:"p"},"<middleware_version_key>")," key in the JSON struct should be replaced by the actual name of the key for the corresponding middleware (e.g. ",(0,i.kt)("inlineCode",{parentName:"p"},"fee_version"),")."),(0,i.kt)("p",null,"During the handshake callbacks, the middleware can unmarshal the version string and retrieve the middleware and application versions. It can do its negotiation logic on ",(0,i.kt)("inlineCode",{parentName:"p"},"<middleware_version_value>"),", and pass the ",(0,i.kt)("inlineCode",{parentName:"p"},"<application_version_value>")," to the underlying application."),(0,i.kt)("p",null,"The middleware should simply pass the capability in the callback arguments along to the underlying application so that it may be claimed by the base application. The base application will then pass the capability up the stack in order to authenticate an outgoing packet/acknowledgement."),(0,i.kt)("p",null,"In the case where the middleware wishes to send a packet or acknowledgment without the involvement of the underlying application, it should be given access to the same ",(0,i.kt)("inlineCode",{parentName:"p"},"scopedKeeper")," as the base application so that it can retrieve the capabilities by itself."),(0,i.kt)("h3",{id:"handshake-callbacks"},"Handshake callbacks"),(0,i.kt)("h4",{id:"onchanopeninit"},(0,i.kt)("inlineCode",{parentName:"h4"},"OnChanOpenInit")),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},'func (im IBCModule) OnChanOpenInit(\n    ctx sdk.Context,\n    order channeltypes.Order,\n    connectionHops []string,\n    portID string,\n    channelID string,\n    channelCap *capabilitytypes.Capability,\n    counterparty channeltypes.Counterparty,\n    version string,\n) (string, error) {\n    if version != "" {\n        // try to unmarshal JSON-encoded version string and pass\n        // the app-specific version to app callback.\n        // otherwise, pass version directly to app callback.\n        metadata, err := Unmarshal(version)\n        if err != nil {\n            // Since it is valid for fee version to not be specified,\n            // the above middleware version may be for another middleware.\n            // Pass the entire version string onto the underlying application.\n            return im.app.OnChanOpenInit(\n                ctx,\n                order,\n                connectionHops,\n                portID,\n                channelID,\n                channelCap,\n                counterparty,\n                version,\n            )\n        }\n    else {\n        metadata = {\n            // set middleware version to default value\n            MiddlewareVersion: defaultMiddlewareVersion,\n            // allow application to return its default version\n            AppVersion: "",\n        }\n    }\n\n    doCustomLogic()\n\n    // if the version string is empty, OnChanOpenInit is expected to return\n    // a default version string representing the version(s) it supports\n    appVersion, err := im.app.OnChanOpenInit(\n        ctx,\n        order,\n        connectionHops,\n        portID,\n        channelID,\n        channelCap,\n        counterparty,\n        metadata.AppVersion, // note we only pass app version here\n    )\n    if err != nil {\n        return "", err\n    }\n\n    version := constructVersion(metadata.MiddlewareVersion, appVersion)\n\n    return version, nil\n}\n')),(0,i.kt)("p",null,"See ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc-go/blob/48a6ae512b4ea42c29fdf6c6f5363f50645591a2/modules/apps/29-fee/ibc_middleware.go#L34-L82"},"here")," an example implementation of this callback for the ICS29 Fee Middleware module."),(0,i.kt)("h4",{id:"onchanopentry"},(0,i.kt)("inlineCode",{parentName:"h4"},"OnChanOpenTry")),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},'func OnChanOpenTry(\n    ctx sdk.Context,\n    order channeltypes.Order,\n    connectionHops []string,\n    portID,\n    channelID string,\n    channelCap *capabilitytypes.Capability,\n    counterparty channeltypes.Counterparty,\n    counterpartyVersion string,\n) (string, error) {\n    // try to unmarshal JSON-encoded version string and pass\n    // the app-specific version to app callback.\n    // otherwise, pass version directly to app callback.\n    cpMetadata, err := Unmarshal(counterpartyVersion)\n    if err != nil {\n        return app.OnChanOpenTry(\n            ctx,\n            order,\n            connectionHops,\n            portID,\n            channelID,\n            channelCap,\n            counterparty,\n            counterpartyVersion,\n        )\n    }\n\n    doCustomLogic()\n\n    // Call the underlying application\'s OnChanOpenTry callback.\n    // The try callback must select the final app-specific version string and return it.\n    appVersion, err := app.OnChanOpenTry(\n        ctx,\n        order,\n        connectionHops,\n        portID,\n        channelID,\n        channelCap,\n        counterparty,\n        cpMetadata.AppVersion, // note we only pass counterparty app version here\n    )\n    if err != nil {\n        return "", err\n    }\n\n    // negotiate final middleware version\n    middlewareVersion := negotiateMiddlewareVersion(cpMetadata.MiddlewareVersion)\n    version := constructVersion(middlewareVersion, appVersion)\n\n    return version, nil\n}\n')),(0,i.kt)("p",null,"See ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc-go/blob/48a6ae512b4ea42c29fdf6c6f5363f50645591a2/modules/apps/29-fee/ibc_middleware.go#L84-L124"},"here")," an example implementation of this callback for the ICS29 Fee Middleware module."),(0,i.kt)("h4",{id:"onchanopenack"},(0,i.kt)("inlineCode",{parentName:"h4"},"OnChanOpenAck")),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},"func OnChanOpenAck(\n    ctx sdk.Context,\n    portID,\n    channelID string,\n    counterpartyChannelID string,\n    counterpartyVersion string,\n) error {\n    // try to unmarshal JSON-encoded version string and pass\n    // the app-specific version to app callback.\n    // otherwise, pass version directly to app callback.\n    cpMetadata, err = UnmarshalJSON(counterpartyVersion)\n    if err != nil {\n        return app.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)\n    }\n\n    if !isCompatible(cpMetadata.MiddlewareVersion) {\n        return error\n    }\n    doCustomLogic()\n\n    // call the underlying application's OnChanOpenTry callback\n    return app.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, cpMetadata.AppVersion)\n}\n")),(0,i.kt)("p",null,"See ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc-go/blob/48a6ae512b4ea42c29fdf6c6f5363f50645591a2/modules/apps/29-fee/ibc_middleware.go#L126-L152"},"here")," an example implementation of this callback for the ICS29 Fee Middleware module."),(0,i.kt)("h3",{id:"onchanopenconfirm"},(0,i.kt)("inlineCode",{parentName:"h3"},"OnChanOpenConfirm")),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},"func OnChanOpenConfirm(\n    ctx sdk.Context,\n    portID,\n    channelID string,\n) error {\n    doCustomLogic()\n\n    return app.OnChanOpenConfirm(ctx, portID, channelID)\n}\n")),(0,i.kt)("p",null,"See ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc-go/blob/48a6ae512b4ea42c29fdf6c6f5363f50645591a2/modules/apps/29-fee/ibc_middleware.go#L154-L162"},"here")," an example implementation of this callback for the ICS29 Fee Middleware module."),(0,i.kt)("h4",{id:"onchancloseinit"},(0,i.kt)("inlineCode",{parentName:"h4"},"OnChanCloseInit")),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},"func OnChanCloseInit(\n    ctx sdk.Context,\n    portID,\n    channelID string,\n) error {\n    doCustomLogic()\n\n    return app.OnChanCloseInit(ctx, portID, channelID)\n}\n")),(0,i.kt)("p",null,"See ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc-go/blob/48a6ae512b4ea42c29fdf6c6f5363f50645591a2/modules/apps/29-fee/ibc_middleware.go#L164-L187"},"here")," an example implementation of this callback for the ICS29 Fee Middleware module."),(0,i.kt)("h4",{id:"onchancloseconfirm"},(0,i.kt)("inlineCode",{parentName:"h4"},"OnChanCloseConfirm")),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},"func OnChanCloseConfirm(\n    ctx sdk.Context,\n    portID,\n    channelID string,\n) error {\n    doCustomLogic()\n\n    return app.OnChanCloseConfirm(ctx, portID, channelID)\n}\n")),(0,i.kt)("p",null,"See ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc-go/blob/48a6ae512b4ea42c29fdf6c6f5363f50645591a2/modules/apps/29-fee/ibc_middleware.go#L189-L212"},"here")," an example implementation of this callback for the ICS29 Fee Middleware module."),(0,i.kt)("p",null,(0,i.kt)("strong",{parentName:"p"},"NOTE"),": Middleware that does not need to negotiate with a counterparty middleware on the remote stack will not implement the version unmarshalling and negotiation, and will simply perform its own custom logic on the callbacks without relying on the counterparty behaving similarly."),(0,i.kt)("h3",{id:"packet-callbacks"},"Packet callbacks"),(0,i.kt)("p",null,"The packet callbacks just like the handshake callbacks wrap the application's packet callbacks. The packet callbacks are where the middleware performs most of its custom logic. The middleware may read the packet flow data and perform some additional packet handling, or it may modify the incoming data before it reaches the underlying application. This enables a wide degree of usecases, as a simple base application like token-transfer can be transformed for a variety of usecases by combining it with custom middleware."),(0,i.kt)("h4",{id:"onrecvpacket"},(0,i.kt)("inlineCode",{parentName:"h4"},"OnRecvPacket")),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},"func OnRecvPacket(\n    ctx sdk.Context,\n    packet channeltypes.Packet,\n    relayer sdk.AccAddress,\n) ibcexported.Acknowledgement {\n    doCustomLogic(packet)\n\n    ack := app.OnRecvPacket(ctx, packet, relayer)\n\n    doCustomLogic(ack) // middleware may modify outgoing ack\n    return ack\n}\n")),(0,i.kt)("p",null,"See ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc-go/blob/48a6ae512b4ea42c29fdf6c6f5363f50645591a2/modules/apps/29-fee/ibc_middleware.go#L214-L237"},"here")," an example implementation of this callback for the ICS29 Fee Middleware module."),(0,i.kt)("h4",{id:"onacknowledgementpacket"},(0,i.kt)("inlineCode",{parentName:"h4"},"OnAcknowledgementPacket")),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},"func OnAcknowledgementPacket(\n    ctx sdk.Context,\n    packet channeltypes.Packet,\n    acknowledgement []byte,\n    relayer sdk.AccAddress,\n) error {\n    doCustomLogic(packet, ack)\n\n    return app.OnAcknowledgementPacket(ctx, packet, ack, relayer)\n}\n")),(0,i.kt)("p",null,"See ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc-go/blob/48a6ae512b4ea42c29fdf6c6f5363f50645591a2/modules/apps/29-fee/ibc_middleware.go#L239-L292"},"here")," an example implementation of this callback for the ICS29 Fee Middleware module."),(0,i.kt)("h4",{id:"ontimeoutpacket"},(0,i.kt)("inlineCode",{parentName:"h4"},"OnTimeoutPacket")),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},"func OnTimeoutPacket(\n    ctx sdk.Context,\n    packet channeltypes.Packet,\n    relayer sdk.AccAddress,\n) error {\n    doCustomLogic(packet)\n\n    return app.OnTimeoutPacket(ctx, packet, relayer)\n}\n")),(0,i.kt)("p",null,"See ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc-go/blob/48a6ae512b4ea42c29fdf6c6f5363f50645591a2/modules/apps/29-fee/ibc_middleware.go#L294-L334"},"here")," an example implementation of this callback for the ICS29 Fee Middleware module."),(0,i.kt)("h3",{id:"ics-4-wrappers"},"ICS-4 wrappers"),(0,i.kt)("p",null,"Middleware must also wrap ICS-4 so that any communication from the application to the ",(0,i.kt)("inlineCode",{parentName:"p"},"channelKeeper")," goes through the middleware first. Similar to the packet callbacks, the middleware may modify outgoing acknowledgements and packets in any way it wishes."),(0,i.kt)("h4",{id:"sendpacket"},(0,i.kt)("inlineCode",{parentName:"h4"},"SendPacket")),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},"func SendPacket(\n    ctx sdk.Context,\n    chanCap *capabilitytypes.Capability,\n    appPacket exported.PacketI,\n) {\n    // middleware may modify packet\n    packet = doCustomLogic(appPacket)\n\n    return ics4Keeper.SendPacket(ctx, chanCap, packet)\n}\n")),(0,i.kt)("p",null,"See ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc-go/blob/48a6ae512b4ea42c29fdf6c6f5363f50645591a2/modules/apps/29-fee/ibc_middleware.go#L336-L343"},"here")," an example implementation of this function for the ICS29 Fee Middleware module."),(0,i.kt)("h4",{id:"writeacknowledgement"},(0,i.kt)("inlineCode",{parentName:"h4"},"WriteAcknowledgement")),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},"// only called for async acks\nfunc WriteAcknowledgement(\n    ctx sdk.Context,\n    chanCap *capabilitytypes.Capability,\n    packet exported.PacketI,\n    ack exported.Acknowledgement,\n) {\n    // middleware may modify acknowledgement\n    ack_bytes = doCustomLogic(ack)\n\n    return ics4Keeper.WriteAcknowledgement(packet, ack_bytes)\n}\n")),(0,i.kt)("p",null,"See ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc-go/blob/48a6ae512b4ea42c29fdf6c6f5363f50645591a2/modules/apps/29-fee/ibc_middleware.go#L345-L353"},"here")," an example implementation of this function for the ICS29 Fee Middleware module."),(0,i.kt)("h4",{id:"getappversion"},(0,i.kt)("inlineCode",{parentName:"h4"},"GetAppVersion")),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-go"},'// middleware must return the underlying application version\nfunc GetAppVersion(\n    ctx sdk.Context,\n    portID,\n    channelID string,\n) (string, bool) {\n    version, found := ics4Keeper.GetAppVersion(ctx, portID, channelID)\n    if !found {\n        return "", false\n    }\n\n    if !MiddlewareEnabled {\n        return version, true\n    }\n\n    // unwrap channel version\n    metadata, err := Unmarshal(version)\n    if err != nil {\n        panic(fmt.Errof("unable to unmarshal version: %w", err))\n    }\n\n    return metadata.AppVersion, true\n}\n')),(0,i.kt)("p",null,"See ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/cosmos/ibc-go/blob/48a6ae512b4ea42c29fdf6c6f5363f50645591a2/modules/apps/29-fee/ibc_middleware.go#L355-L358"},"here")," an example implementation of this function for the ICS29 Fee Middleware module."))}m.isMDXComponent=!0}}]);