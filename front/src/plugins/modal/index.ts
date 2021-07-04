import { App, createApp, inject, watch } from "vue";
import EmptyModal from "./components/Empty.vue";
import InfoModal from "./components/Info.vue";
import ErrorModal from "./components/Error.vue";

export enum ModalType {
    Empty = 1,
    Info,
    Error
}

function getModal(type: ModalType){
    switch (type) {
        case ModalType.Info:
            return InfoModal
        case ModalType.Error:
            return ErrorModal
        default:
            return EmptyModal
    }
}

export function showModal(type: ModalType, title: string, body:string, time=2000, closeCallback=()=>{}){
    const mountPoint = document.createElement('div');
    mountPoint.className = "w-full h-full fixed t-0 l-0"
    document.body.appendChild(mountPoint);

    const modal = createApp(getModal(type),{
        title: title,
        body: body,
        time: time,
        close:()=>{document.body.removeChild(mountPoint);closeCallback();}
    }).mount(mountPoint);

    return modal
}

export function showError(title: string, body:string, time=2000, closeCallback=()=>{}){
    return showModal(ModalType.Error,title,body,time, closeCallback)
}

export function showInfo(title: string, body:string, time=2000, closeCallback=()=>{}){
    return showModal(ModalType.Info,title,body,time, closeCallback)
}

export function showEmpty(title: string, body:string, time=2000, closeCallback=()=>{}){
    return showModal(ModalType.Empty,title,body,time, closeCallback)
}
