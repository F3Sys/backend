import{u as z,f as H,r as f,d as M,g as x,h as N,c as y,a as e,t as m,e as r,w as k,b as R,i as G,j as J,v as O,k as _,l as I,F,m as L,o as g}from"./index-DDo2_VMk.js";import{C as q,A as K,p as Q,a as X,b as U,L as W,P as $,d as Y,e as Z,S as ee,u as te,f as ae,g as se,h as ne,D as oe,i as le,F as re,_ as ie}from"./_plugin-vue_export-helper-DRxNFOdN.js";const S=z(),ue=H("entry",()=>{const w=f(0),c=f([]),n=f([]),u=f([]);function p(){w.value=0,c.value=[],n.value=[],u.value=[]}async function T(){const[a]=await Promise.all([b(),C(),v(),j()]);return a}async function h(a){const o=new Headers;o.append("Authorization","Bearer "+S.key),o.append("Content-Type","application/json");const s="https://api.aicj.io/";try{return(await fetch(s+"protected/push/entry",{method:"POST",headers:o,body:JSON.stringify({f3sid:a})})).ok}catch{return!1}}async function v(){const a=new Headers;a.append("Authorization","Bearer "+S.key),a.append("Content-Type","application/json");const o="https://api.aicj.io/";try{const s=await fetch(o+"protected/table",{method:"GET",headers:a}).then(d=>d.json());return n.value=s,!0}catch{return!1}}async function b(){const a=new Headers;a.append("Authorization","Bearer "+S.key),a.append("Content-Type","application/json");const o="https://api.aicj.io/";try{const s=await fetch(o+"protected/count",{method:"GET",headers:a}).then(d=>d.json());return w.value=s.count,!0}catch{return!1}}async function C(){const a=new Headers;a.append("Authorization","Bearer "+S.key),a.append("Content-Type","application/json");const o="https://api.aicj.io/";try{const s=await fetch(o+"protected/entry_count",{method:"GET",headers:a}).then(d=>d.json());return c.value=s,!0}catch{return!1}}async function j(){const a=new Headers;a.append("Authorization","Bearer "+S.key),a.append("Content-Type","application/json");const o="https://api.aicj.io/";try{const s=await fetch(o+"protected/data/entry",{method:"GET",headers:a}).then(d=>d.json());return u.value=s,!0}catch{return!1}}return{counts:c,table:n,line_graph_data:u,count:w,sendEntry:h,getTable:v,getCount:b,getEntryCount:C,getData:j,update:T,clear:p}}),de={class:"container-fluid"},ce={class:"parent"},pe={class:"Order"},he={class:"mb-0"},fe=["disabled"],ve=["aria-invalid"],me={key:0,id:"f3sid-helper"},ye=["aria-busy","disabled"],ge={class:"flex items-center mx-auto"},be={class:"Stats"},_e={class:"text-nowrap"},Se={class:"striped mb-0"},we={scope:"row"},Ce={class:"Table"},je={class:"text-nowrap"},Te={class:"striped mb-0"},ke={scope:"row"},xe={class:"Donut-Chart"},Fe={style:{zoom:"1.1"}},Le={class:"Line-Chart"},ze={style:{zoom:"1.1"}},Ee=["open"],De={class:"max-w-lg"},Ae=21,Ve=M({__name:"EntryView",setup(w){q.register(K,Q,X,U,W,$,Y,Z);const c=z(),n=ue(),u=f(!1),p=f(!1),T=f(!1),h=f(""),v=f(!1),b=new ee({minLength:7,alphabet:"23456789CFGHJMPQRVWX",blocklist:new Set([])}),C=async([i])=>{b.decode(i.rawValue).length===2&&(h.value=i.rawValue,T.value=!0,p.value=!1)},j=()=>{p.value=!p.value},a=async()=>{if(u.value=!0,h.value===""||b.decode(h.value).length!==2?(u.value=!1,v.value=!0):v.value=!1,v.value){u.value=!1;return}await n.sendEntry(h.value)&&await n.update()&&(h.value=""),u.value=!1},o=x(()=>n.counts.map(i=>i.type)),s=x(()=>n.counts.map(i=>i.count));let d=Array(Ae).fill(0).reduce((i,t,l)=>(l%2===0?i.push((l/2+8).toString().padStart(2,"0")+":00"):i.push((Math.floor(l/2)+8).toString().padStart(2,"0")+":30"),i),[]);const{charging:E,chargingTime:D,dischargingTime:A,level:V}=te(),{resume:B,isActive:P}=ae(()=>{c.sendStatus(E.value,D.value,A.value,V.value)},6e4,{immediate:!1});return N(()=>{(async()=>await n.update())(),!P.value&&c.canSendStatus&&B()}),(i,t)=>(g(),y("main",de,[e("div",ce,[e("div",pe,[e("article",null,[e("header",null,[e("hgroup",he,[e("h2",null,m(r(c).name),1),e("p",null,m(r(c).type.charAt(0)+r(c).type.toLowerCase().slice(1)),1)])]),e("form",{onSubmit:k(a,["prevent"])},[e("fieldset",{disabled:u.value},[e("label",null,[t[3]||(t[3]=R(" F3SiD ")),e("fieldset",{onKeydown:t[1]||(t[1]=G(k(()=>{},["prevent"]),["enter"])),role:"group",style:{"margin-top":"calc(var(--pico-spacing) * .25)"},"aria-invalid":v.value,"aria-describedby":"f3sid-helper"},[J(e("input",{"onUpdate:modelValue":t[0]||(t[0]=l=>h.value=l),name:"f3sid",class:"!h-10"},null,512),[[O,h.value]]),e("button",{type:"submit",onClick:k(j,["prevent"]),class:"flex items-center justify-center !size-10 !p-0"},[_(r(se),{class:"m-2 size-6"})])],40,ve),v.value?(g(),y("small",me," F3SiDが無効です ")):I("",!0)])],8,fe),e("button",{type:"submit",class:"mb-0 flex items-center","aria-busy":u.value,disabled:u.value},[e("div",ge,[_(r(ne),{class:"size-4 mr-1"}),t[4]||(t[4]=e("span",{class:"text-size-lg"},"送信する",-1))])],8,ye)],32)])]),e("div",be,[e("article",_e,[t[7]||(t[7]=e("header",null,[e("hgroup",{class:"mb-0"},[e("h2",null,"統計"),e("p",null,"Stats")])],-1)),e("table",Se,[t[6]||(t[6]=e("thead",null,[e("tr",null,[e("th",{scope:"col"},"種類"),e("th",{scope:"col"},"回数")])],-1)),e("tbody",null,[(g(!0),y(F,null,L(r(n).counts,l=>(g(),y("tr",null,[e("th",we,m(l.type),1),e("td",null,m(l.count.toLocaleString("ja-JP")),1)]))),256))]),e("tfoot",null,[e("tr",null,[t[5]||(t[5]=e("th",{scope:"row",class:"text-align-center"},"合計",-1)),e("td",null,m(r(n).count.toLocaleString("ja-JP")),1)])])])])]),e("div",Ce,[e("article",je,[t[9]||(t[9]=e("header",null,[e("hgroup",{class:"mb-0"},[e("h2",null,"表"),e("p",null,"Table")])],-1)),e("table",Te,[t[8]||(t[8]=e("thead",null,[e("tr",null,[e("th",{scope:"col"},"F3SiD"),e("th",{scope:"col"},"種類"),e("th",{scope:"col"},"時刻")])],-1)),e("tbody",null,[(g(!0),y(F,null,L(r(n).table,l=>(g(),y("tr",null,[e("th",ke,m(l.f3sid),1),e("td",null,m(l.type),1),e("td",null,m(new Date(l.created_at).toLocaleTimeString("ja-JP")),1)]))),256))])])])]),e("div",xe,[e("article",Fe,[t[10]||(t[10]=e("header",null,[e("hgroup",{class:"mb-0"},[e("h2",null,"種別人数割合"),e("p",null,"Percentage of total count by type")])],-1)),_(r(oe),{data:{labels:o.value,datasets:[{label:"回数",data:s.value}]},options:{animation:!1,plugins:{colors:{forceOverride:!0},legend:{labels:{font:{size:15}}}}}},null,8,["data"])])]),e("div",Le,[e("article",ze,[t[11]||(t[11]=e("header",null,[e("hgroup",{class:"mb-0"},[e("h2",null,"時間における回数のグラフ"),e("p",null,"The graph of count vs time")])],-1)),_(r(le),{data:{labels:r(d),datasets:r(n).line_graph_data},options:{animation:!1,plugins:{colors:{forceOverride:!0},legend:{labels:{font:{size:15}}}},scales:{x:{ticks:{font:{size:12}}},y:{ticks:{font:{size:12}}}}}},null,8,["data"])])])]),e("dialog",{open:p.value},[e("article",De,[t[12]||(t[12]=e("header",null,[e("hgroup",{style:{"margin-bottom":"0px"}},[e("h2",null,"F3SiD"),e("p",null,"Scan the F3SiD")])],-1)),_(r(re),{paused:!p.value,onDetect:C},null,8,["paused"]),e("footer",null,[e("button",{onClick:t[2]||(t[2]=l=>{p.value=!p.value}),class:"secondary"}," Close ")])])],8,Ee)]))}}),He=ie(Ve,[["__scopeId","data-v-51408727"]]);export{He as default};
