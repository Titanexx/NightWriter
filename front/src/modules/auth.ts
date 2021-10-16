import { ref } from 'vue';
import { useApi, useApiToken } from "./api";
import { Crypto } from './crypto';
import router from '../routes/index';

var _token: string = "";
const user= ref();
const isAuthenticated = ref(false);

// console.log(useApi);

function setToken(token: string){
    const jwt = JSON.parse(atob(token.split('.')[1]));
    const timeout = new Date(jwt.exp * 1000).getTime() - Date.now() - (30 * 1000);
    setTimeout(() => {logout()}, timeout); // handle refresh URL
    var { data, get } = useApi("/api/users/me",token);
    get().then(() => {
        user.value = data.value;
        if(!Crypto.decryptPem({public:user.value.public_key,private:user.value.private_key})){
            
            router.push({ name: 'decrypt-key' });
        }
    })
    
    window.localStorage.setItem("token", token)
    _token = token;
    isAuthenticated.value = true;
}

function loadKeys(masterpassword: string, remember=false){
    const pems = {
        public: user.value.public_key,
        private: user.value.private_key
    }
    return Crypto.loadKeys(pems, masterpassword, remember)
}

function getToken(){
    if (isValidToken()){
        return _token;
    } else {
        return undefined;
    }
}

function logout(){
    useApiToken("/api/logout").get()
    router.push({ name: 'login' })
    localStorage.removeItem('token');
    _token = "";
    isAuthenticated.value = false;
    user.value = undefined;
    Crypto.clearLocalStorage();
}

function isValidToken(){
	return isAuthenticated.value && _token.length > 0;
}

function isAuth(){
    return isAuthenticated.value;
}

const token = window.localStorage.getItem("token")
if ( token ) {
    setToken(token);
}

export const Auth = {
    user,
    isAuthenticated,
    isAuth,
    loadKeys,
    setToken,
    getToken,
    logout
}
