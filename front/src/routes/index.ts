import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import { Auth } from '../modules/auth'

const routes = [
    {
        path: '/',
        name: 'home',
        component: Home,
        meta: { requiresAuth: true },
    },
    {
        path: '/register',
        name: 'register',
        // route level code-splitting
        // this generates a separate chunk (about.[hash].js) for this route
        // which is lazy-loaded when the route is visited.
        component: () => import(/* webpackChunkName: "auth" */ '../views/Register.vue'),
        meta: { requiresAuth: false },
    },
    {
        path: '/login',
        name: 'login',
        component: () => import(/* webpackChunkName: "auth" */ '../views/Login.vue'),
        meta: { requiresAuth: false },
    },
    {
        path: '/decrypt-key',
        name: 'decrypt-key',
        component: () => import(/* webpackChunkName: "auth" */ '../views/DecryptKey.vue'),
        meta: { requiresAuth: true },
    },
    {
        path: '/users/',
        name: 'users',
        component: () => import(/* webpackChunkName: "user" */ '../views/Users/List.vue'),
        meta: { requiresAuth: true },
    },
    {
        path: '/users/:id',
        name: 'users-id',
        component: () => import(/* webpackChunkName: "user" */ '../views/Users/User.vue'),
        meta: { requiresAuth: true },
    },
    {
        path: '/docs/',
        name: 'docs',
        component: () => import(/* webpackChunkName: "doc" */ '../views/Docs/List.vue'),
        meta: { requiresAuth: true },
    },
    {
        path: '/docs/:id',
        name: 'docs-id',
        component: () => import(/* webpackChunkName: "doc" */ '../views/Docs/Doc.vue'),
        meta: { requiresAuth: true },
    },
    {
        path: '/docs/new',
        name: 'docs-new',
        component: () => import(/* webpackChunkName: "doc" */ '../views/Docs/Add.vue'),
        meta: { requiresAuth: true },
    },
    {
        path: '/*',
        redirect: '/',
    },
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    //@ts-ignore
    routes
})

router.beforeEach((to, from, next) => {
    // Not logged into a guarded route?
    if (to.meta.requiresAuth === true && !Auth.isAuth()) {
        next({ name: 'login' })
    } else next()
})

export default router
