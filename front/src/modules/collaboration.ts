import * as Y from 'yjs';

import { ySyncPlugin } from 'y-prosemirror';
import { EditorPlugin } from '@toast-ui/editor/types/editor';

var collabWebSockets: {[url: number]: CollabWS} = {}

export function getCollabWebSocket(id: number){
    return collabWebSockets[id]
}

export class CollabWS {
    id: number
    url: string
    doc: Y.Doc
    type: Y.XmlFragment = undefined
    ws: WebSocket
    encrypt: Function
    decrypt: Function
    initCallback: Function
    isFirstClient: boolean
    isInitialized: boolean
    mustSendUpdate: boolean
    
    constructor(url: string,  encrypt: Function, decrypt: Function, initFinishedCallback: Function){
        this.id = Math.floor(Math.random() * 100000)
        while(this.id in collabWebSockets){
            this.id = Math.floor(Math.random() * 100000)
        }
        this.url = url
        this.isFirstClient = false
        this.isInitialized = false
        this.mustSendUpdate = true

        this.encrypt = encrypt
        this.decrypt = decrypt
        this.initCallback = initFinishedCallback
        
        this.doc = new Y.Doc()
        this.doc.on('update',this.updateYDocHandler(this))
        this.type = this.doc.getXmlFragment('prosemirror')

        collabWebSockets[this.id] = this
        console.log("cws.constructor")

    }

    updateYDocHandler(self){
        return (update) => {
            if(import.meta.env.DEV)
                console.log("cws.updateYDocHandler", self.mustSendUpdate)
            if(self.mustSendUpdate){
                // console.log("cws.updateYDocHandler")
                self.ws.send(self.encrypt(update))
            }
        }
    }
    
    initFinished(){
        console.log("cws.initFinished")
        this.isInitialized = true
    }

    initFirstClient(){
        if(import.meta.env.DEV)
            console.log("cws.initFirstClient")
        this.isFirstClient= true
        this.initCallback()
    }
    
    initOtherClient1stPart(){
        if(import.meta.env.DEV)
            console.log("cws.initOtherClient1stPart")
        this.isFirstClient= false
    }

    initOtherClient2ndPart(data){
        if(import.meta.env.DEV)
            console.log("cws.initOtherClient2ndPart")
        this.mustSendUpdate = false
        data.forEach((u)=>{
            Y.applyUpdate(this.doc, this.decrypt(u))
        })
        this.mustSendUpdate = true
        this.initCallback()
    }

    editorIsLoaded(content: string){
        if(import.meta.env.DEV)
            console.log("cws.editorIsLoaded")
        if(this.isFirstClient && content == ""){
            this.ws.send(this.encrypt(Y.encodeStateAsUpdate(this.doc)))
        }
    }

    readUpdate(data){
        if(import.meta.env.DEV)
            console.log("cws.readUpdate")
        this.mustSendUpdate = false
        Y.applyUpdate(this.doc, this.decrypt(data))
        this.mustSendUpdate = true
    }

    init(){
        var self = this
        this.ws = new WebSocket(this.url)
        this.ws.onmessage = function(event) {
            var data = JSON.parse(event.data)

            if (data == 0) {// Init is finished
                self.initFinished()
            } else if (data == 1){
                self.initFirstClient()
            } else if (data == 2) {
                self.initOtherClient1stPart()
            } else {
                if(typeof(data) == "object"){ // If it's an object (array is object), it's init part
                    self.initOtherClient2ndPart(data)
                }else{ // else it's an classic yjs update
                    self.readUpdate(data)
                }
            }
        }
    }

    get plugins(): EditorPlugin[]{
        return [
            (c)=>{return {markdownPlugins: [(c)=>{ return ySyncPlugin(this.type)}]}},
        ]
    }

    close(){
        this.ws.close()
        this.doc.destroy()
        delete collabWebSockets[this.id]
    }
}