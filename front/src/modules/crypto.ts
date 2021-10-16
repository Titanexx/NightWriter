import { showError } from '@plugins/modal';
import * as forge from 'node-forge';
import { reactive } from 'vue';
import { toBlobArray } from './helpers';

const keyIts = 100000
const rsaBits = 2048 // 2048, 4096
const alg = 'aes256'
const keyLen = 32
const ivLen = 12

interface Keys {
    public: any,
    private: any,
}

var keys = reactive<Keys>({
    public: undefined,
    private: undefined
})

function genKeys(password: string, next: Function, ifErr: Function){
    const salt = forge.random.getBytesSync(256);
    const key = forge.pkcs5.pbkdf2(password, salt, keyIts, keyLen);
    forge.pki.rsa.generateKeyPair({bits: rsaBits, e: 0x10001}, (err,keypair) => {
        if(err === null){
            const pem = {
                public: forge.pki.publicKeyToPem(keypair.publicKey),
                private: forge.util.encode64(salt) + "|" + forge.pki.encryptRsaPrivateKey(keypair.privateKey, key,{algorithm: alg}),
            }
            next(pem)
        } else {
            console.log("Error in key generation:",err)
            ifErr()
        }
    });
}

function genIV(){
    return forge.util.bytesToHex(forge.random.getBytesSync(ivLen));
}

function loadKeys(pem: Keys, password:string, remember=false) {
    var private_data = pem.private.split('|');
    const salt = forge.util.decode64(private_data[0]);
    const key = forge.pkcs5.pbkdf2(password, salt, keyIts, keyLen);
    
    //@ts-ignore
    const privateKey = forge.pki.decryptRsaPrivateKey(private_data[1],key,{algorithm: alg}) 

    if(privateKey == null){
        keys.public = undefined
        keys.private = undefined
        return false
    } else {
        
        
        keys.public = forge.pki.publicKeyFromPem(pem.public)
        keys.private = privateKey

        if(remember){
            window.localStorage.setItem("masterpassword",password)
        }

        return true
    }
}

function encryptRSA(data: string, publicKey?: string){
    var key;
    if(publicKey === undefined){
        key = keys.public;
    } else {
        key= forge.pki.publicKeyFromPem(publicKey)
    }
    return forge.util.encode64(key.encrypt(data))
}

function decryptRSA(data: string){
    return keys.private.decrypt(forge.util.decode64(data))
}

function encryptAES(data: string | Uint8Array, key?: string, iv?:string): {data: string, key: string, iv: string}{
    if(key === undefined) {
        key = forge.random.getBytesSync(keyLen);
    } else {
        key = forge.util.hexToBytes(key);
    }

    if(iv === undefined) {
        iv = forge.random.getBytesSync(ivLen);
    } else {
        iv = forge.util.hexToBytes(iv);
    }

    let dataBuffer = forge.util.createBuffer(data )
    let cipher = forge.cipher.createCipher('AES-GCM', key);
    cipher.start({iv: iv});
    cipher.update(dataBuffer);
    cipher.finish();
    const encrypted = forge.util.encode64(cipher.output.getBytes());
    const tag = cipher.mode.tag.toHex();

    return {
        data: tag+encrypted,
        key: forge.util.bytesToHex(key),
        iv: forge.util.bytesToHex(iv),
    }
}

function decryptAES(data: string, key: string, iv: string): {data: string, iv: string}{
    if(data === undefined){
        return {
            data: "",
            iv: forge.util.bytesToHex(iv),
        }
    }
    key = forge.util.hexToBytes(key)
    iv = forge.util.hexToBytes(iv)
    const tag = forge.util.createBuffer(forge.util.hexToBytes(data.slice(0, 32)))
    data = forge.util.decode64(data.slice(32))

    var decipher = forge.cipher.createDecipher('AES-GCM', key);
    decipher.start({iv: iv, tag: tag});
    decipher.update(forge.util.createBuffer(data));
    const pass = decipher.finish();
    
    if(!pass) {
        showError("Decryption error","The decryption can't be done. There is a problem with data.")  
        return undefined
    }

    return {
        data: decipher.output.bytes(),
        iv: forge.util.bytesToHex(iv),
    }
}

function encryptBlobAES(data: Blob | File, key: string): any{
    key = forge.util.hexToBytes(key);
    var iv = forge.random.getBytesSync(ivLen);
    var cipher = forge.cipher.createCipher('AES-GCM', key);
    cipher.start({iv:iv})
    //@ts-ignore
    var promise = new Promise<any>((resolve, reject) => {
        data.arrayBuffer().then((buffer) => {
            cipher.update(forge.util.createBuffer(buffer));
            cipher.finish();
            
            const encrypted = cipher.output.getBytes();
            const tag = cipher.mode.tag.getBytes();

            resolve(new Blob([iv+tag+encrypted], {type: 'application/octet-stream'}))
        })
    })

    return promise
}

function decryptBlobAES(data, key){
    key = forge.util.hexToBytes(key);
    var iv = forge.util.createBuffer(data.slice(0,12))
    var tag = forge.util.createBuffer(data.slice(12,28))

    data = forge.util.createBuffer(data.slice(28))
    

    var decipher = forge.cipher.createDecipher('AES-GCM', key);
    decipher.start({iv: iv, tag: tag});
    decipher.update(data);
    const pass = decipher.finish();

    if(!pass) {
        showError("Decryption error","The decryption can't be done. There is a problem with data.")  
        return undefined
    }
    var res = decipher.output.getBytes()
 
    return new Blob([toBlobArray(res)],{type:"image/jpeg"})
}

function decryptPem(pem: Keys):boolean{
    const masterpassword = window.localStorage.getItem("masterpassword")
    if (masterpassword != undefined && loadKeys(pem,masterpassword,true)){
        console.log("Your private key is fully decrypted with your password stored in localstorage")
        return true
    } else {
        return false
    }
}

function clearLocalStorage(){
    window.localStorage.removeItem("masterpassword")
}

export const Crypto = {
    genKeys,genIV,
    loadKeys,
    decryptPem: decryptPem,
    clearLocalStorage,
    encryptAES, decryptAES,
    encryptRSA, decryptRSA,
    encryptBlobAES, decryptBlobAES,
}