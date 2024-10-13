import{u as z,g as H,r as f,d as M,h as x,i as N,c as y,a as e,t as m,e as l,w as k,b as R,j as G,k as J,v as O,f as _,l as I,F,m as L,o as g}from"./index-2LCq3qsk.js";import{F as q}from"./vue-qrcode-reader-DJ1kwrmJ.js";import{C as K,A as Q,p as X,a as U,b as W,L as $,P as Y,d as Z,e as ee,S as te,u as ae,f as se,g as ne,h as oe,D as le,i as re,_ as ie}from"./_plugin-vue_export-helper-C59vZBZt.js";const S=z(),ue=H("entry",()=>{const w=f(0),d=f([]),n=f([]),i=f([]);function p(){w.value=0,d.value=[],n.value=[],i.value=[]}async function T(){const[a]=await Promise.all([b(),C(),v(),j()]);return a}async function h(a){const o=new Headers;o.append("Authorization","Bearer "+S.key),o.append("Content-Type","application/json");const u="https://api.aicj.io/";try{return(await fetch(u+"protected/push/entry",{method:"POST",headers:o,body:JSON.stringify({f3sid:a})})).ok}catch{return!1}}async function v(){const a=new Headers;a.append("Authorization","Bearer "+S.key),a.append("Content-Type","application/json");const o="https://api.aicj.io/";try{const u=await fetch(o+"protected/table",{method:"GET",headers:a}).then(c=>c.json());return n.value=u,!0}catch{return!1}}async function b(){const a=new Headers;a.append("Authorization","Bearer "+S.key),a.append("Content-Type","application/json");const o="https://api.aicj.io/";try{const u=await fetch(o+"protected/count",{method:"GET",headers:a}).then(c=>c.json());return w.value=u.count,!0}catch{return!1}}async function C(){const a=new Headers;a.append("Authorization","Bearer "+S.key),a.append("Content-Type","application/json");const o="https://api.aicj.io/";try{const u=await fetch(o+"protected/entry_count",{method:"GET",headers:a}).then(c=>c.json());return d.value=u,!0}catch{return!1}}async function j(){const a=new Headers;a.append("Authorization","Bearer "+S.key),a.append("Content-Type","application/json");const o="https://api.aicj.io/";try{const u=await fetch(o+"protected/data/entry",{method:"GET",headers:a}).then(c=>c.json());return i.value=u,!0}catch{return!1}}return{counts:d,table:n,line_graph_data:i,count:w,sendEntry:h,getTable:v,getCount:b,getEntryCount:C,getData:j,update:T,clear:p}}),ce={class:"container-fluid"},de={class:"parent"},pe={class:"Order"},he={class:"mb-0"},fe=["disabled"],ve=["aria-invalid"],me={key:0,id:"f3sid-helper"},ye=["aria-busy","disabled"],ge={class:"flex items-center mx-auto"},be={class:"Stats"},_e={class:"text-nowrap"},Se={class:"striped mb-0"},we={scope:"row"},Ce={class:"Table"},je={class:"text-nowrap"},Te={class:"striped mb-0"},ke={scope:"row"},xe={class:"Donut-Chart"},Fe={style:{zoom:"1.1"}},Le={class:"Line-Chart"},ze={style:{zoom:"1.1"}},Ee=["open"],De={class:"max-w-lg"},Ae=21,Ve=M({__name:"EntryView",setup(w){K.register(Q,X,U,W,$,Y,Z,ee);const d=z(),n=ue(),i=f(!1),p=f(!1),T=f(!1),h=f(""),v=f(!1),b=new te({minLength:7,alphabet:"23456789CFGHJMPQRVWX",blocklist:new Set([])}),C=async([r])=>{b.decode(r.rawValue).length===2&&(h.value=r.rawValue,T.value=!0,p.value=!1)},j=()=>{p.value=!p.value},a=async()=>{if(i.value=!0,h.value===""||b.decode(h.value).length!==2?(i.value=!1,v.value=!0):v.value=!1,v.value){i.value=!1;return}await n.sendEntry(h.value)&&await n.update()&&(h.value=""),i.value=!1},o=x(()=>n.counts.map(r=>r.type)),u=x(()=>n.counts.map(r=>r.count)),c=Array(Ae).fill(0).reduce((r,t,s)=>(s%2===0?r.push((s/2+8).toString().padStart(2,"0")+":00"):r.push((Math.floor(s/2)+8).toString().padStart(2,"0")+":30"),r),[]),{charging:E,chargingTime:D,dischargingTime:A,level:V}=ae(),{resume:B,isActive:P}=se(async()=>{await d.sendStatus(E.value,D.value,A.value,V.value)},6e4,{immediate:!1});return N(()=>{(async()=>await n.update())(),!P.value&&d.canSendStatus&&B()}),(r,t)=>(g(),y("main",ce,[e("div",de,[e("div",pe,[e("article",null,[e("header",null,[e("hgroup",he,[e("h2",null,m(l(d).name),1),e("p",null,m(l(d).type.charAt(0)+l(d).type.toLowerCase().slice(1)),1)])]),e("form",{onSubmit:k(a,["prevent"])},[e("fieldset",{disabled:i.value},[e("label",null,[t[3]||(t[3]=R(" F3SiD ")),e("fieldset",{onKeydown:t[1]||(t[1]=G(k(()=>{},["prevent"]),["enter"])),role:"group",style:{"margin-top":"calc(var(--pico-spacing) * .25)"},"aria-invalid":v.value,"aria-describedby":"f3sid-helper"},[J(e("input",{"onUpdate:modelValue":t[0]||(t[0]=s=>h.value=s),name:"f3sid",class:"!h-10"},null,512),[[O,h.value]]),e("button",{type:"submit",onClick:k(j,["prevent"]),class:"flex items-center justify-center !size-10 !p-0"},[_(l(ne),{class:"m-2 size-6"})])],40,ve),v.value?(g(),y("small",me," F3SiDが無効です ")):I("",!0)])],8,fe),e("button",{type:"submit",class:"mb-0 flex items-center","aria-busy":i.value,disabled:i.value},[e("div",ge,[_(l(oe),{class:"size-4 mr-1"}),t[4]||(t[4]=e("span",{class:"text-size-lg"},"送信する",-1))])],8,ye)],32)])]),e("div",be,[e("article",_e,[t[7]||(t[7]=e("header",null,[e("hgroup",{class:"mb-0"},[e("h2",null,"統計"),e("p",null,"Stats")])],-1)),e("table",Se,[t[6]||(t[6]=e("thead",null,[e("tr",null,[e("th",{scope:"col"},"種類"),e("th",{scope:"col"},"回数")])],-1)),e("tbody",null,[(g(!0),y(F,null,L(l(n).counts,s=>(g(),y("tr",{key:s.type},[e("th",we,m(s.type),1),e("td",null,m(s.count.toLocaleString("ja-JP")),1)]))),128))]),e("tfoot",null,[e("tr",null,[t[5]||(t[5]=e("th",{scope:"row",class:"text-align-center"},"合計",-1)),e("td",null,m(l(n).count.toLocaleString("ja-JP")),1)])])])])]),e("div",Ce,[e("article",je,[t[9]||(t[9]=e("header",null,[e("hgroup",{class:"mb-0"},[e("h2",null,"表"),e("p",null,"Table")])],-1)),e("table",Te,[t[8]||(t[8]=e("thead",null,[e("tr",null,[e("th",{scope:"col"},"F3SiD"),e("th",{scope:"col"},"種類"),e("th",{scope:"col"},"時刻")])],-1)),e("tbody",null,[(g(!0),y(F,null,L(l(n).table,s=>(g(),y("tr",{key:s.f3sid+s.created_at},[e("th",ke,m(s.f3sid),1),e("td",null,m(s.type),1),e("td",null,m(new Date(s.created_at).toLocaleTimeString("ja-JP")),1)]))),128))])])])]),e("div",xe,[e("article",Fe,[t[10]||(t[10]=e("header",null,[e("hgroup",{class:"mb-0"},[e("h2",null,"種別人数割合"),e("p",null,"Percentage of total count by type")])],-1)),_(l(le),{data:{labels:o.value,datasets:[{label:"回数",data:u.value}]},options:{animation:!1,plugins:{colors:{forceOverride:!0},legend:{labels:{font:{size:15}}}}}},null,8,["data"])])]),e("div",Le,[e("article",ze,[t[11]||(t[11]=e("header",null,[e("hgroup",{class:"mb-0"},[e("h2",null,"時間における回数のグラフ"),e("p",null,"The graph of count vs time")])],-1)),_(l(re),{data:{labels:l(c),datasets:l(n).line_graph_data},options:{animation:!1,plugins:{colors:{forceOverride:!0},legend:{labels:{font:{size:15}}}},scales:{x:{ticks:{font:{size:12}}},y:{ticks:{font:{size:12}}}}}},null,8,["data"])])])]),e("dialog",{open:p.value},[e("article",De,[t[12]||(t[12]=e("header",null,[e("hgroup",{style:{"margin-bottom":"0px"}},[e("h2",null,"F3SiD"),e("p",null,"Scan the F3SiD")])],-1)),_(l(q),{paused:!p.value,onDetect:C},null,8,["paused"]),e("footer",null,[e("button",{onClick:t[2]||(t[2]=s=>{p.value=!p.value}),class:"secondary"}," Close ")])])],8,Ee)]))}}),Me=ie(Ve,[["__scopeId","data-v-a2c46d9c"]]);export{Me as default};
