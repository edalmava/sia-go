import{C as e,D as t,S as n,a as r,c as i,d as a,f as o,g as s,h as c,i as l,l as u,m as d,o as f,p,s as m,t as h,u as g,x as _,y as v}from"../ui-dNetwNdp.js";/* empty css                   */t((()=>{s(),v(),m();var t=[],y=[],b=null,x=[];document.addEventListener(`DOMContentLoaded`,()=>{if(g()){if(!a([`usuarios_`,`roles_`,`permisos_`])){u(`No tienes permisos para acceder a esta página`,`error`),setTimeout(()=>{window.location.href=`dashboard.html`},1500);return}c(),d(),p(),f(),w(),C(),T(),E(),O(),M(),S()}});async function S(){try{x=(await n.getAll()).data||[]}catch(e){console.error(`Error preloading permissions:`,e)}}function C(){let e=document.getElementById(`addRoleBtn`);e&&(e.style.display=o(`roles_crear`)?`inline-flex`:`none`)}function w(){let e=document.querySelectorAll(`.tab-btn`),t=document.querySelectorAll(`.tab-content`),n=o(`roles_ver`)||o(`roles_crear`)||o(`roles_editar`)||o(`roles_eliminar`),r=o(`permisos_ver`);e.forEach(i=>{let a=i.dataset.tab;if(a===`roles`&&!n){i.style.display=`none`;return}if(a===`permisos`&&!r){i.style.display=`none`;return}i.addEventListener(`click`,()=>{e.forEach(e=>e.classList.remove(`active`)),t.forEach(e=>e.classList.remove(`active`)),i.classList.add(`active`),document.getElementById(`${i.dataset.tab}-tab`)?.classList.add(`active`)})});let i=Array.from(e).filter(e=>e.style.display!==`none`);if(i.length>0){i[0].classList.add(`active`);let e=i[0].dataset.tab;document.getElementById(`${e}-tab`)?.classList.add(`active`)}}function T(){document.getElementById(`addRoleBtn`)?.addEventListener(`click`,()=>openRoleModal()),document.getElementById(`closeRoleModal`)?.addEventListener(`click`,F),document.getElementById(`cancelRoleBtn`)?.addEventListener(`click`,F),document.getElementById(`roleModal`)?.addEventListener(`click`,e=>{e.target.id===`roleModal`&&F()}),document.getElementById(`roleForm`)?.addEventListener(`submit`,I),document.getElementById(`permisoSearch`)?.addEventListener(`input`,j),document.getElementById(`moduloFilter`)?.addEventListener(`change`,j)}async function E(){let t=document.getElementById(`rolesTableBody`);t.innerHTML=`<tr><td colspan="6" class="loading-cell">Cargando roles...</td></tr>`;try{D((await e.getAll()).data||[])}catch(e){console.error(`Error:`,e),t.innerHTML=`<tr><td colspan="6" class="empty-cell">Error al cargar roles</td></tr>`,r(e,`Error al cargar roles`)}}function D(e){let t=document.getElementById(`rolesTableBody`),n=o(`roles_editar`),r=o(`roles_eliminar`);if(e.length===0){t.innerHTML=`<tr><td colspan="6" class="empty-cell">No hay roles registrados</td></tr>`;return}t.innerHTML=e.map(e=>{let t=``;return n&&(t+=`
                <button class="btn btn-icon btn-ghost" onclick="openRoleModal(${e.id_rol}, '${l(e.nombre)}', '${l(e.descripcion||``)}', ${e.es_rol_sistema})" title="Editar rol">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" width="18" height="18">
                        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                    </svg>
                </button>`),r&&!e.es_rol_sistema&&(t+=`
                <button class="btn btn-icon btn-ghost" style="color: var(--error);" onclick="deleteRole(${e.id_rol}, '${l(e.nombre)}')" title="Eliminar rol">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" width="18" height="18">
                        <polyline points="3 6 5 6 21 6"/>
                        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                    </svg>
                </button>`),`
            <tr>
                <td>${e.id_rol}</td>
                <td><strong>${l(e.nombre)}</strong></td>
                <td>${l(e.descripcion||`-`)}</td>
                <td><span class="badge badge-${e.es_rol_sistema?`primary`:`secondary`}">${e.es_rol_sistema?`Sí`:`No`}</span></td>
                <td><span class="badge badge-info">${e.permisos_count||0} permisos</span></td>
                <td>
                    <div class="action-buttons">
                        ${t}
                    </div>
                </td>
            </tr>
        `}).join(``)}async function O(){let e=document.getElementById(`permisosTableBody`);e.innerHTML=`<tr><td colspan="5" class="loading-cell">Cargando permisos...</td></tr>`;try{t=(await n.getAll()).data||[],await k(),A(t)}catch(t){console.error(`Error:`,t),e.innerHTML=`<tr><td colspan="5" class="empty-cell">Error al cargar permisos</td></tr>`,r(t,`Error al cargar permisos`)}}async function k(){try{y=(await _.getAll()).data||[];let e=document.getElementById(`moduloFilter`);y.forEach(t=>{let n=document.createElement(`option`);n.value=t.id_modulo,n.textContent=t.nombre,e.appendChild(n)})}catch(e){console.error(`Error loading modulos:`,e)}}function A(e){let t=document.getElementById(`permisosTableBody`);if(e.length===0){t.innerHTML=`<tr><td colspan="5" class="empty-cell">No hay permisos registrados</td></tr>`;return}t.innerHTML=e.map(e=>`
        <tr>
            <td>${e.id_permiso}</td>
            <td><strong>${l(e.nombre)}</strong></td>
            <td><code>${l(e.codigo)}</code></td>
            <td>${l(e.modulo_nombre||`-`)}</td>
            <td>${l(e.descripcion||`-`)}</td>
        </tr>
    `).join(``)}function j(){let e=document.getElementById(`permisoSearch`).value.toLowerCase(),n=document.getElementById(`moduloFilter`).value,r=t;e&&(r=r.filter(t=>t.nombre.toLowerCase().includes(e)||t.codigo.toLowerCase().includes(e)||t.descripcion&&t.descripcion.toLowerCase().includes(e))),n&&(r=r.filter(e=>e.id_modulo==n)),A(r)}async function M(){let e=document.getElementById(`modulosGrid`);e.innerHTML=`<div class="loading-card">Cargando módulos...</div>`;try{N((await _.getAll()).data||[])}catch(t){console.error(`Error:`,t),e.innerHTML=`<div class="empty-state">Error al cargar módulos</div>`,r(t,`Error al cargar módulos`)}}function N(e){let t=document.getElementById(`modulosGrid`);if(e.length===0){t.innerHTML=`<div class="empty-state">No hay módulos registrados</div>`;return}t.innerHTML=e.map(e=>`
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
                    <div class="font-semibold text-primary">${l(e.nombre)}</div>
                    <div class="text-xs text-muted">${l(e.codigo)}</div>
                </div>
            </div>
            <p class="text-sm text-secondary mb-4">${l(e.descripcion||`Sin descripción`)}</p>
            <div class="flex items-center gap-2 text-xs text-muted">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" width="14" height="14">
                    <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                    <polyline points="22 4 12 14.01 9 11.01"/>
                </svg>
                <span>${e.permisos_count||0} permisos</span>
            </div>
        </div>
    `).join(``)}window.openRoleModal=async function(t=null,n=``,a=``,o=!1){let s=document.getElementById(`roleModal`),c=document.getElementById(`roleForm`),l=document.getElementById(`roleModalTitle`),u=document.getElementById(`rolePermissions`);if(b=t,c.reset(),u.innerHTML=`<div class="text-sm text-muted p-4">Cargando permisos...</div>`,x.length===0&&await S(),t){l.textContent=`Editar Rol`,document.getElementById(`roleId`).value=t,document.getElementById(`roleName`).value=n,document.getElementById(`roleDesc`).value=a,document.getElementById(`roleSistema`).checked=o,document.getElementById(`roleName`).disabled=o;try{P(((await e.getById(t)).data?.permisos||[]).map(e=>e.id_permiso))}catch(e){r(e,`Error al cargar permisos del rol`)}}else l.textContent=`Nuevo Rol`,document.getElementById(`roleId`).value=``,document.getElementById(`roleName`).disabled=!1,P([]);i(s)};function P(e=[]){let t=document.getElementById(`rolePermissions`);if(x.length===0){t.innerHTML=`<div class="text-error p-4">No se cargaron los permisos correctamente</div>`;return}let n={};x.forEach(e=>{n[e.id_modulo]||(n[e.id_modulo]={nombre:e.modulo_nombre,permisos:[]}),n[e.id_modulo].permisos.push(e)});let r=``;Object.entries(n).forEach(([,t])=>{r+=`<div class="permission-module-header">${l(t.nombre)}</div>`,t.permisos.forEach(t=>{let n=e.includes(t.id_permiso);r+=`
                <label class="checkbox-group mb-2">
                    <input type="checkbox" name="permisos" class="checkbox-input" value="${t.id_permiso}" ${n?`checked`:``}>
                    <span class="text-sm">${l(t.nombre)}</span>
                </label>
            `})}),t.innerHTML=r}function F(){h(document.getElementById(`roleModal`)),b=null}async function I(t){if(t.preventDefault(),!o(`roles_crear`)&&!o(`roles_editar`)){u(`No tienes permisos para guardar roles`,`error`);return}let n=t.target,i=new FormData(n),a={nombre:i.get(`nombre`),descripcion:i.get(`descripcion`),es_rol_sistema:i.get(`es_rol_sistema`)===`on`},s=n.querySelectorAll(`input[name="permisos"]:checked`),c=Array.from(s).map(e=>parseInt(e.value));try{b?(await e.update(b,{...a,permisos:c}),u(`Rol actualizado exitosamente`,`success`)):(await e.create({...a,permisos:c}),u(`Rol creado exitosamente`,`success`)),F(),E()}catch(e){r(e,`Error al guardar rol`)}}window.deleteRole=async function(t,n){if(!o(`roles_eliminar`)){u(`No tienes permisos para eliminar roles`,`error`);return}if(confirm(`¿Está seguro de eliminar el rol "${n}"?`))try{await e.delete(t),u(`Rol eliminado exitosamente`,`success`),E()}catch(e){r(e,`Error al eliminar rol`)}}}))();