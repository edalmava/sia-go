import{E as e,S as t,a as n,b as r,c as i,d as a,f as o,h as s,i as c,l,m as u,o as d,p as f,s as p,t as m,u as h,v as g,x as _}from"../ui-CtHqQ12e.js";/* empty css                   */e((()=>{s(),g(),d();var e=[],v=[],y=null,b=[];document.addEventListener(`DOMContentLoaded`,()=>{if(l()){if(!h([`usuarios_`,`roles_`,`permisos_`])){i(`No tienes permisos para acceder a esta página`,`error`),setTimeout(()=>{window.location.href=`dashboard.html`},1500);return}u(),f(),o(),C(),S(),w(),T(),D(),j(),x()}});async function x(){try{b=(await _.getAll()).data||[]}catch(e){console.error(`Error preloading permissions:`,e)}}function S(){let e=document.getElementById(`addRoleBtn`);e&&(e.style.display=a(`roles_crear`)?`inline-flex`:`none`)}function C(){let e=document.querySelectorAll(`.tab-btn`),t=document.querySelectorAll(`.tab-content`),n=a(`roles_ver`)||a(`roles_crear`)||a(`roles_editar`)||a(`roles_eliminar`),r=a(`permisos_ver`);e.forEach(i=>{let a=i.dataset.tab;if(a===`roles`&&!n){i.style.display=`none`;return}if(a===`permisos`&&!r){i.style.display=`none`;return}i.addEventListener(`click`,()=>{e.forEach(e=>e.classList.remove(`active`)),t.forEach(e=>e.classList.remove(`active`)),i.classList.add(`active`),document.getElementById(`${i.dataset.tab}-tab`)?.classList.add(`active`)})});let i=Array.from(e).filter(e=>e.style.display!==`none`);if(i.length>0){i[0].classList.add(`active`);let e=i[0].dataset.tab;document.getElementById(`${e}-tab`)?.classList.add(`active`)}}function w(){document.getElementById(`addRoleBtn`)?.addEventListener(`click`,()=>openRoleModal()),document.getElementById(`closeRoleModal`)?.addEventListener(`click`,P),document.getElementById(`cancelRoleBtn`)?.addEventListener(`click`,P),document.getElementById(`roleModal`)?.addEventListener(`click`,e=>{e.target.id===`roleModal`&&P()}),document.getElementById(`roleForm`)?.addEventListener(`submit`,F),document.getElementById(`permisoSearch`)?.addEventListener(`input`,A),document.getElementById(`moduloFilter`)?.addEventListener(`change`,A)}async function T(){let e=document.getElementById(`rolesTableBody`);e.innerHTML=`<tr><td colspan="6" class="loading-cell">Cargando roles...</td></tr>`;try{E((await t.getAll()).data||[])}catch(t){console.error(`Error:`,t),e.innerHTML=`<tr><td colspan="6" class="empty-cell">Error al cargar roles</td></tr>`,n(t,`Error al cargar roles`)}}function E(e){let t=document.getElementById(`rolesTableBody`),n=a(`roles_editar`),r=a(`roles_eliminar`);if(e.length===0){t.innerHTML=`<tr><td colspan="6" class="empty-cell">No hay roles registrados</td></tr>`;return}t.innerHTML=e.map(e=>{let t=``;return n&&(t+=`
                <button class="btn btn-icon btn-ghost" onclick="openRoleModal(${e.id_rol}, '${c(e.nombre)}', '${c(e.descripcion||``)}', ${e.es_rol_sistema})" title="Editar rol">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" width="18" height="18">
                        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                    </svg>
                </button>`),r&&!e.es_rol_sistema&&(t+=`
                <button class="btn btn-icon btn-ghost" style="color: var(--error);" onclick="deleteRole(${e.id_rol}, '${c(e.nombre)}')" title="Eliminar rol">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" width="18" height="18">
                        <polyline points="3 6 5 6 21 6"/>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                    </svg>
                </button>`),`
            <tr>
                <td>${e.id_rol}</td>
                <td><strong>${c(e.nombre)}</strong></td>
                <td>${c(e.descripcion||`-`)}</td>
                <td><span class="badge badge-${e.es_rol_sistema?`primary`:`secondary`}">${e.es_rol_sistema?`Sí`:`No`}</span></td>
                <td><span class="badge badge-info">${e.permisos_count||0} permisos</span></td>
                <td>
                    <div class="action-buttons">
                        ${t}
                    </div>
                </td>
            </tr>
        `}).join(``)}async function D(){let t=document.getElementById(`permisosTableBody`);t.innerHTML=`<tr><td colspan="5" class="loading-cell">Cargando permisos...</td></tr>`;try{e=(await _.getAll()).data||[],await O(),k(e)}catch(e){console.error(`Error:`,e),t.innerHTML=`<tr><td colspan="5" class="empty-cell">Error al cargar permisos</td></tr>`,n(e,`Error al cargar permisos`)}}async function O(){try{v=(await r.getAll()).data||[];let e=document.getElementById(`moduloFilter`);v.forEach(t=>{let n=document.createElement(`option`);n.value=t.id_modulo,n.textContent=t.nombre,e.appendChild(n)})}catch(e){console.error(`Error loading modulos:`,e)}}function k(e){let t=document.getElementById(`permisosTableBody`);if(e.length===0){t.innerHTML=`<tr><td colspan="5" class="empty-cell">No hay permisos registrados</td></tr>`;return}t.innerHTML=e.map(e=>`
        <tr>
            <td>${e.id_permiso}</td>
            <td><strong>${c(e.nombre)}</strong></td>
            <td><code>${c(e.codigo)}</code></td>
            <td>${c(e.modulo_nombre||`-`)}</td>
            <td>${c(e.descripcion||`-`)}</td>
        </tr>
    `).join(``)}function A(){let t=document.getElementById(`permisoSearch`).value.toLowerCase(),n=document.getElementById(`moduloFilter`).value,r=e;t&&(r=r.filter(e=>e.nombre.toLowerCase().includes(t)||e.codigo.toLowerCase().includes(t)||e.descripcion&&e.descripcion.toLowerCase().includes(t))),n&&(r=r.filter(e=>e.id_modulo==n)),k(r)}async function j(){let e=document.getElementById(`modulosGrid`);e.innerHTML=`<div class="loading-card">Cargando módulos...</div>`;try{M((await r.getAll()).data||[])}catch(t){console.error(`Error:`,t),e.innerHTML=`<div class="empty-state">Error al cargar módulos</div>`,n(t,`Error al cargar módulos`)}}function M(e){let t=document.getElementById(`modulosGrid`);if(e.length===0){t.innerHTML=`<div class="empty-state">No hay módulos registrados</div>`;return}t.innerHTML=e.map(e=>`
        <div class="card">
            <div class="flex items-center gap-4 mb-4">
                <div class="stat-icon institutions" style="width: 40px; height: 40px;">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" width="20" height="20">
                        <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                        <line x1="3" y1="9" x2="21" y2="9"/>
                        <line x1="9" y1="21" x2="9" y2="9"/>
                    </svg>
                </div>
                <div>
                    <div class="font-semibold text-primary">${c(e.nombre)}</div>
                    <div class="text-xs text-muted">${c(e.codigo)}</div>
                </div>
            </div>
            <p class="text-sm text-secondary mb-4">${c(e.descripcion||`Sin descripción`)}</p>
            <div class="flex items-center gap-2 text-xs text-muted">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" width="14" height="14">
                    <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                    <polyline points="22 4 12 14.01 9 11.01"/>
                </svg>
                <span>${e.permisos_count||0} permisos</span>
            </div>
        </div>
    `).join(``)}window.openRoleModal=async function(e=null,r=``,i=``,a=!1){let o=document.getElementById(`roleModal`),s=document.getElementById(`roleForm`),c=document.getElementById(`roleModalTitle`),l=document.getElementById(`rolePermissions`);if(y=e,s.reset(),l.innerHTML=`<div class="text-sm text-muted p-4">Cargando permisos...</div>`,b.length===0&&await x(),e){c.textContent=`Editar Rol`,document.getElementById(`roleId`).value=e,document.getElementById(`roleName`).value=r,document.getElementById(`roleDesc`).value=i,document.getElementById(`roleSistema`).checked=a,document.getElementById(`roleName`).disabled=a;try{N(((await t.getById(e)).data?.permisos||[]).map(e=>e.id_permiso))}catch(e){n(e,`Error al cargar permisos del rol`)}}else c.textContent=`Nuevo Rol`,document.getElementById(`roleId`).value=``,document.getElementById(`roleName`).disabled=!1,N([]);p(o)};function N(e=[]){let t=document.getElementById(`rolePermissions`);if(b.length===0){t.innerHTML=`<div class="text-error p-4">No se cargaron los permisos correctamente</div>`;return}let n={};b.forEach(e=>{n[e.id_modulo]||(n[e.id_modulo]={nombre:e.modulo_nombre,permisos:[]}),n[e.id_modulo].permisos.push(e)});let r=``;Object.entries(n).forEach(([,t])=>{r+=`<div class="permission-module-header">${c(t.nombre)}</div>`,t.permisos.forEach(t=>{let n=e.includes(t.id_permiso);r+=`
                <label class="checkbox-group mb-2">
                    <input type="checkbox" name="permisos" class="checkbox-input" value="${t.id_permiso}" ${n?`checked`:``}>
                    <span class="text-sm">${c(t.nombre)}</span>
                </label>
            `})}),t.innerHTML=r}function P(){m(document.getElementById(`roleModal`)),y=null}async function F(e){if(e.preventDefault(),!a(`roles_crear`)&&!a(`roles_editar`)){i(`No tienes permisos para guardar roles`,`error`);return}let r=e.target,o=new FormData(r),s={nombre:o.get(`nombre`),descripcion:o.get(`descripcion`),es_rol_sistema:o.get(`es_rol_sistema`)===`on`},c=r.querySelectorAll(`input[name="permisos"]:checked`),l=Array.from(c).map(e=>parseInt(e.value));try{y?(await t.update(y,{...s,permisos:l}),i(`Rol actualizado exitosamente`,`success`)):(await t.create({...s,permisos:l}),i(`Rol creado exitosamente`,`success`)),P(),T()}catch(e){n(e,`Error al guardar rol`)}}window.deleteRole=async function(e,r){if(!a(`roles_eliminar`)){i(`No tienes permisos para eliminar roles`,`error`);return}if(confirm(`¿Está seguro de eliminar el rol "${r}"?`))try{await t.delete(e),i(`Rol eliminado exitosamente`,`success`),T()}catch(e){n(e,`Error al eliminar rol`)}}}))();