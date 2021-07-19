import { computed, ref } from 'vue'

export function usePropsWrapper(props:any, emit: any, name = 'modelValue') { 
  return computed({ 
    get: () => props[name], 
    set: (value) => emit(`update:${name}`, value)
  })
}

export function getLast(array: any[]){
    return array[array.length - 1]
}

export function getWSServerUrl(){
    return import.meta.env.VITE_HTTPS == "true" ? `wss://${window.location.host}` : `ws://${window.location.host}`
}

export function findAttachmentIDs(md: string, url: string){
    var attachmentUrl = `\(${url}/([0-9]+)\)`;
    var attachmentRe = new RegExp(attachmentUrl,"g");
    return [...md.matchAll(attachmentRe)].map(e => parseInt(e[2]))
}

export var themeIsLight = ref(!document.documentElement.classList.contains("dark"))

export interface IListOptions{
    val: number | string,
    str: string,
}
export type ListOptionsMap = Record<string | number | symbol, IListOptions>;

export function toBlobArray(str) {
    var res = new Uint8Array(str.length);
    for (var i = 0; i < str.length; i++) {
        res[i] = str.charCodeAt(i);
    }
    return res;
}

/////////////////////////
// Code from tui.editor
const XMLSPECIAL = '[&<>"]';
const reXmlSpecial = new RegExp(XMLSPECIAL, 'g');
function replaceUnsafeChar(s: string) {
    switch (s) {
        case '&':
            return '&amp;';
        case '<':
            return '&lt;';
        case '>':
            return '&gt;';
        case '"':
            return '&quot;';
        default:
            return s;
    }
}

export function escapeXml(s: string) {
    if (reXmlSpecial.test(s)) {
        return s.replace(reXmlSpecial, replaceUnsafeChar);
    }
    return s;
}
/////////////////////////