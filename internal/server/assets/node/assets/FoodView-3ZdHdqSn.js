import{u as O,g as se,r as p,d as le,h as D,i as oe,c as u,a as e,t as h,e as y,w as E,b as J,j as ie,k as T,v as R,f as _,l as I,F as S,m as L,o as r,n as ue,p as re}from"./index-2LCq3qsk.js";import{F as ce}from"./vue-qrcode-reader-DJ1kwrmJ.js";import{c as N,C as de,A as pe,p as he,a as ye,b as ve,L as me,P as fe,d as be,e as ge,S as _e,u as Se,f as ke,g as we,h as Ce,D as Fe,i as je,_ as qe}from"./_plugin-vue_export-helper-C59vZBZt.js";/**
 * @license lucide-vue-next v0.445.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const ze=N("CheckIcon",[["path",{d:"M20 6 9 17l-5-5",key:"1gmf2c"}]]);/**
 * @license lucide-vue-next v0.445.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const xe=N("PencilIcon",[["path",{d:"M21.174 6.812a1 1 0 0 0-3.986-3.987L3.842 16.174a2 2 0 0 0-.5.83l-1.321 4.352a.5.5 0 0 0 .623.622l4.353-1.32a2 2 0 0 0 .83-.497z",key:"1a8usu"}],["path",{d:"m15 5 4 4",key:"1mk7zo"}]]);/**
 * @license lucide-vue-next v0.445.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const Te=N("XIcon",[["path",{d:"M18 6 6 18",key:"1bl5f8"}],["path",{d:"m6 6 12 12",key:"d8bk6v"}]]),k=O(),Le=se("food",()=>{const P=p([]),q=p(0),z=p(0),m=p([]),l=p([]),v=p([]);function f(){P.value=[],z.value=0,q.value=0,m.value=[],l.value=[],v.value=[]}async function M(){const[n]=await Promise.all([g(),b(),F(),C(),A()]);return n}async function x(){const n=new Headers;n.append("Authorization","Bearer "+k.key),n.append("Content-Type","application/json");const o="https://api.aicj.io/";try{const c=await fetch(o+"protected/foods",{method:"GET",headers:n}).then(d=>d.json());return P.value=c,!0}catch{return!1}}async function w(n,o){const c=new Headers;c.append("Authorization","Bearer "+k.key),c.append("Content-Type","application/json");const d="https://api.aicj.io/";try{return(await fetch(d+"protected/push/foodstall",{method:"POST",headers:c,body:JSON.stringify({f3sid:n,foods:o})})).ok}catch{return!1}}async function C(){const n=new Headers;n.append("Authorization","Bearer "+k.key),n.append("Content-Type","application/json");const o="https://api.aicj.io/";try{const c=await fetch(o+"protected/table",{method:"GET",headers:n}).then(d=>d.json());return l.value=c,!0}catch{return!1}}async function b(){const n=new Headers;n.append("Authorization","Bearer "+k.key),n.append("Content-Type","application/json");const o="https://api.aicj.io/";try{const c=await fetch(o+"protected/count",{method:"GET",headers:n}).then(d=>d.json());return z.value=c.count,!0}catch{return!1}}async function g(){const n=new Headers;n.append("Authorization","Bearer "+k.key),n.append("Content-Type","application/json");const o="https://api.aicj.io/";try{const c=await fetch(o+"protected/quantity",{method:"GET",headers:n}).then(d=>d.json());return q.value=c.quantity,!0}catch{return!1}}async function F(){const n=new Headers;n.append("Authorization","Bearer "+k.key),n.append("Content-Type","application/json");const o="https://api.aicj.io/";try{const c=await fetch(o+"protected/food_count",{method:"GET",headers:n}).then(d=>d.json());return m.value=c,!0}catch{return!1}}async function j(n,o,c){const d=new Headers;d.append("Authorization","Bearer "+k.key),d.append("Content-Type","application/json");const V="https://api.aicj.io/";try{return(await fetch(V+"protected/push/foodstall",{method:"PATCH",headers:d,body:JSON.stringify({id:n,food_id:o,quantity:c})})).ok}catch{return!1}}async function A(){const n=new Headers;n.append("Authorization","Bearer "+k.key);const o="https://api.aicj.io/";try{const c=await fetch(o+"protected/data/foodstall",{method:"GET",headers:n}).then(d=>d.json());return v.value=c,!0}catch{return!1}}return{foods:P,quantity:q,counts:m,count:z,table:l,line_graph_data:v,getFoods:x,sendFood:w,updateFood:j,update:M,clear:f}}),Pe={class:"container-fluid"},Ae={class:"parent"},Ve={class:"Order"},De={class:"mb-0"},Me=["disabled"],Be=["aria-invalid"],Ee={key:0,id:"f3sid-helper"},Je=["aria-invalid"],Re=["id","name","onUpdate:modelValue"],Ie=["for"],Ne={key:0,id:"food-helper"},He={key:0},Oe=["onUpdate:modelValue"],Ge=["value"],Ue=["aria-busy","disabled"],Qe={class:"flex items-center mx-auto"},Ye={class:"Stats"},Xe={class:"text-nowrap"},$e={class:"striped mb-0"},Ke={scope:"row",class:"text-sm"},We={class:"Table"},Ze={class:"text-nowrap"},et={class:"striped mb-0"},tt={scope:"row"},at={key:0},nt=["disabled"],st=["value"],lt={key:1},ot={key:2},it=["disabled"],ut={key:3},rt={class:"grid grid-cols-1 gap-2"},ct=["onClick","disabled"],dt=["onClick","disabled"],pt=["onClick","disabled"],ht={class:"Donut-Chart"},yt={style:{zoom:"1.1"}},vt={class:"Line-Chart"},mt={style:{zoom:"1.1"}},ft=["open"],bt={class:"max-w-lg"},gt=21,_t=le({__name:"FoodView",setup(P){de.register(pe,he,ye,ve,me,fe,be,ge);const q=["rgb(54, 162, 235)","rgb(255, 99, 132)","rgb(255, 159, 64)","rgb(255, 205, 86)","rgb(75, 192, 192)","rgb(153, 102, 255)","rgb(201, 203, 207)"],z=p(""),m=O(),l=Le(),v=p(!1),f=p(!1),M=p(!1),x=p(0),w=p(0),C=p(1),b=p(!1),g=p(""),F=p(!1),j=p(!1),A=new _e({minLength:7,alphabet:"23456789CFGHJMPQRVWX",blocklist:new Set([])}),n=p(new Map(l.table.map(s=>[s.id,!1]))),o=p(new Map),c=D(()=>Array.from(o.value.entries()).filter(([,s])=>s.isSelected).reduce((s,[a,t])=>{const i=l.foods.find(ne=>ne.id===a);return s+(i?i.price*t.quantity:0)},0).toLocaleString("ja-JP",{style:"currency",currency:"JPY"})),d=async([s])=>{A.decode(s.rawValue).length===2&&(g.value=s.rawValue,M.value=!0,f.value=!1)},V=()=>{f.value=!f.value},H=async s=>{console.log(x.value,w.value,C.value),b.value=!0,await l.updateFood(Number(x.value),Number(w.value),Number(C.value))&&await l.update()&&(b.value=!1,n.value.set(s,!1))},G=s=>{n.value.set(s,!1)},U=s=>{for(const[a,t]of n.value)if(t&&a!==s.id)return;w.value=s.food_id,C.value=s.quantity,x.value=s.id,n.value.set(s.id,!0)},Q=async()=>{if(v.value=!0,g.value===""||A.decode(g.value).length!==2?(v.value=!1,F.value=!0):F.value=!1,Array.from(o.value.values()).some(t=>t.isSelected&&t.quantity>=1)?j.value=!1:(v.value=!1,j.value=!0),F.value||j.value){v.value=!1;return}await l.sendFood(g.value,Array.from(o.value.entries()).filter(([,t])=>t.isSelected).map(([t,i])=>({id:Number(t),quantity:Number(i.quantity)})))&&await l.update()&&(g.value="",o.value=new Map(l.foods.map(i=>[i.id,{name:i.name,isSelected:!1,quantity:1}]))),v.value=!1},Y=D(()=>l.counts.map(s=>s.name)),B=D(()=>l.counts.map(s=>s.count)),X=D(()=>l.counts.map(s=>s.quantity)),$=Array(gt).fill(0).reduce((s,a,t)=>(t%2===0?s.push((t/2+8).toString().padStart(2,"0")+":00"):s.push((Math.floor(t/2)+8).toString().padStart(2,"0")+":30"),s),[]),{charging:K,chargingTime:W,dischargingTime:Z,level:ee}=Se(),{resume:te,isActive:ae}=ke(async()=>{await m.sendStatus(K.value,W.value,Z.value,ee.value)},6e4,{immediate:!1});return oe(()=>{(async()=>(await l.getFoods(),await l.update(),z.value=l.counts.reduce((s,a)=>s+a.price*a.quantity,0).toLocaleString("ja-JP",{style:"currency",currency:"JPY"}),o.value=new Map(l.foods.map(s=>[s.id,{name:s.name,isSelected:!1,quantity:1}]))))(),!ae.value&&m.canSendStatus&&te()}),(s,a)=>(r(),u("main",Pe,[e("div",Ae,[e("div",Ve,[e("article",null,[e("header",null,[e("hgroup",De,[e("h2",null,h(y(m).name),1),e("p",null,h(y(m).type.charAt(0)+y(m).type.toLowerCase().slice(1)),1)])]),e("form",{onSubmit:E(Q,["prevent"])},[e("fieldset",{disabled:v.value},[e("label",null,[a[5]||(a[5]=J(" F3SiD ")),e("fieldset",{onKeydown:a[1]||(a[1]=ie(E(()=>{},["prevent"]),["enter"])),role:"group",style:{"margin-top":"calc(var(--pico-spacing) * .25)"},"aria-invalid":F.value,"aria-describedby":"f3sid-helper"},[T(e("input",{"onUpdate:modelValue":a[0]||(a[0]=t=>g.value=t),name:"f3sid",class:"!h-10"},null,512),[[R,g.value]]),e("button",{type:"submit",onClick:E(V,["prevent"]),class:"flex items-center justify-center !size-10 !p-0"},[_(y(we),{class:"m-2 size-6"})])],40,Be),F.value?(r(),u("small",Ee," F3SiDが無効です ")):I("",!0)]),e("fieldset",{"aria-invalid":j.value,"aria-describedby":"food-helper"},[a[6]||(a[6]=e("legend",null,"商品:",-1)),(r(!0),u(S,null,L(o.value,t=>(r(),u(S,{key:t[0].toString()},[T(e("input",{type:"checkbox",id:t[0].toString(),name:t[1].name,"onUpdate:modelValue":i=>t[1].isSelected=i},null,8,Re),[[ue,t[1].isSelected]]),e("label",{for:t[0].toString()},h(t[1].name),9,Ie)],64))),128))],8,Je),j.value?(r(),u("small",Ne," 商品を選択してください ")):I("",!0),(r(!0),u(S,null,L(o.value,t=>(r(),u(S,{key:t[0]},[t[1].isSelected?(r(),u("label",He,[J(h(t[1].name)+" ",1),T(e("input",{type:"number",name:"quantity",min:"1","onUpdate:modelValue":i=>t[1].quantity=i,class:"!h-10"},null,8,Oe),[[R,t[1].quantity]])])):I("",!0)],64))),128)),e("label",null,[a[7]||(a[7]=J(" 値段 ")),e("input",{value:c.value,name:"price",class:"!h-10",readonly:""},null,8,Ge)])],8,Me),e("button",{type:"submit",class:"mb-0 flex items-center","aria-busy":v.value,disabled:v.value},[e("div",Qe,[_(y(Ce),{class:"size-4 mr-1"}),a[8]||(a[8]=e("span",{class:"text-size-lg"},"送信する",-1))])],8,Ue)],32)])]),e("div",Ye,[e("article",Xe,[a[11]||(a[11]=e("header",null,[e("hgroup",{class:"mb-0"},[e("h2",null,"統計"),e("p",null,"Stats")])],-1)),e("div",null,[e("table",$e,[a[10]||(a[10]=e("thead",null,[e("tr",null,[e("th",{scope:"col"},"商品名"),e("th",{scope:"col"},"値段"),e("th",{scope:"col"},"回数"),e("th",{scope:"col"},"個数"),e("th",{scope:"col"},"総額")])],-1)),e("tbody",null,[(r(!0),u(S,null,L(y(l).counts,t=>(r(),u("tr",{key:t.id},[e("th",Ke,h(t.name),1),e("td",null,h(t.price.toLocaleString("ja-JP",{style:"currency",currency:"JPY"})),1),e("td",null,h(t.count),1),e("td",null,h(t.quantity),1),e("td",null,h((t.quantity*t.price).toLocaleString("ja-JP",{style:"currency",currency:"JPY"})),1)]))),128))]),e("tfoot",null,[e("tr",null,[a[9]||(a[9]=e("th",{scope:"row",colspan:"2",class:"text-align-center"},"合計",-1)),e("td",null,h(y(l).count),1),e("td",null,h(y(l).quantity),1),e("td",null,h(z.value),1)])])])])])]),e("div",We,[e("article",Ze,[a[13]||(a[13]=e("header",null,[e("hgroup",{class:"mb-0"},[e("h2",null,"表"),e("p",null,"Table")])],-1)),e("table",et,[a[12]||(a[12]=e("thead",null,[e("tr",null,[e("th",{scope:"col"},"F3SiD"),e("th",{scope:"col"},"商品名"),e("th",{scope:"col"},"回数"),e("th",{scope:"col"},"値段"),e("th",{scope:"col"},"購入時刻"),e("th",{scope:"col",class:"text-center"},"編集")])],-1)),e("tbody",null,[(r(!0),u(S,null,L(y(l).table,t=>(r(),u("tr",{key:t.f3sid+t.created_at},[e("th",tt,h(t.f3sid),1),n.value.get(t.id)?(r(),u("td",at,[T(e("select",{disabled:b.value,name:"newFoodId","onUpdate:modelValue":a[2]||(a[2]=i=>w.value=i),class:"!mb-0 !py-0 !h-8",required:""},[(r(!0),u(S,null,L(y(l).foods,i=>(r(),u("option",{value:i.id,key:i.id},h(i.name),9,st))),128))],8,nt),[[re,w.value]])])):(r(),u("td",lt,h(t.food_name),1)),n.value.get(t.id)?(r(),u("td",ot,[T(e("input",{type:"number","onUpdate:modelValue":a[3]||(a[3]=i=>C.value=i),class:"!mb-0 !py-0 !h-8 !w-20",disabled:b.value},null,8,it),[[R,C.value]])])):(r(),u("td",ut,h(t.quantity),1)),e("td",null,h((t.price*t.quantity).toLocaleString("ja-JP",{style:"currency",currency:"JPY"})),1),e("td",null,h(new Date(t.created_at).toLocaleTimeString("ja-JP")),1),e("td",null,[e("div",rt,[n.value.get(t.id)?(r(),u(S,{key:0},[e("button",{onClick:i=>G(t.id),disabled:b.value,class:"flex items-center justify-center size-8 p-0 mr-auto outline secondary"},[_(y(Te),{class:"m-1 size-5"})],8,ct),e("button",{onClick:i=>H(t.id),disabled:b.value,class:"flex items-center justify-center size-8 p-0 outline"},[_(y(ze),{class:"m-1 size-5"})],8,dt)],64)):(r(),u("button",{key:1,onClick:i=>U(t),disabled:b.value,class:"flex items-center justify-center size-8 p-0 mx-auto outline contrast"},[_(y(xe),{class:"m-2 size-3"})],8,pt))])])]))),128))])])])]),e("div",ht,[e("article",yt,[a[14]||(a[14]=e("header",null,[e("hgroup",{class:"mb-0"},[e("h2",null,"セット別売上割合"),e("p",null,"Percentage of total count and quantity by set")])],-1)),_(y(Fe),{data:{labels:Y.value,datasets:[{label:"回数",data:B.value,backgroundColor:[...q.slice(0,B.value.length)]},{label:"個数",data:X.value,backgroundColor:[...q.slice(0,B.value.length)]}]},options:{animation:!1,plugins:{legend:{labels:{font:{size:15}}}}}},null,8,["data"])])]),e("div",vt,[e("article",mt,[a[15]||(a[15]=e("header",null,[e("hgroup",{class:"mb-0"},[e("h2",null,"時間における回数と個数のグラフ"),e("p",null,"The graph of count and quantity vs time")])],-1)),_(y(je),{data:{labels:y($),datasets:y(l).line_graph_data},options:{animation:!1,plugins:{colors:{forceOverride:!0},legend:{labels:{font:{size:15}}}},scales:{x:{ticks:{font:{size:12}}},y:{ticks:{font:{size:12}}}}}},null,8,["data"])])])]),e("dialog",{open:f.value},[e("article",bt,[a[16]||(a[16]=e("header",null,[e("hgroup",{style:{"margin-bottom":"0px"}},[e("h2",null,"F3SiD"),e("p",null,"Scan the F3SiD")])],-1)),_(y(ce),{paused:!f.value,onDetect:d},null,8,["paused"]),e("footer",null,[e("button",{onClick:a[4]||(a[4]=t=>{f.value=!f.value}),class:"secondary"}," Close ")])])],8,ft)]))}}),Ct=qe(_t,[["__scopeId","data-v-9c5695b9"]]);export{Ct as default};
