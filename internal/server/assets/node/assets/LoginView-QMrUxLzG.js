import{d,u as r,r as s,c,a,w as b,b as v,e as m,v as p,t as f,o as _}from"./index-DY-KYBUN.js";const w={class:"container max-w-lg"},g=["aria-busy","disabled"],y=["aria-busy","disabled"],N=d({__name:"LoginView",setup(x){const n=r(),o=s(""),e=s(!1),l=s("Login");async function i(){l.value="Loading...",e.value=!e.value,await n.sendOTP(o.value)&&await n.getNode(),window.location.reload()}return(S,t)=>(_(),c("main",w,[a("article",null,[a("form",{onSubmit:b(i,["prevent"])},[a("fieldset",null,[a("label",null,[t[1]||(t[1]=v(" OTP ")),m(a("input",{"onUpdate:modelValue":t[0]||(t[0]=u=>o.value=u),"aria-busy":e.value,disabled:e.value,name:"otp"},null,8,g),[[p,o.value]])])]),a("button",{type:"submit","aria-busy":e.value,disabled:e.value},f(l.value),9,y)],32)])]))}});export{N as default};
