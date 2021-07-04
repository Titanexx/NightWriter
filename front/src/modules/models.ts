import JSZip from "jszip"
import { saveAs } from 'file-saver';
import { useApiToken } from "./api"
import { Crypto } from "./crypto"
import { getWSServerUrl, findAttachmentIDs, ListOptionsMap } from "./helpers"

abstract class Searchable {
    abstract search(text:string): boolean;
}

export enum DocRights {
    NoRight = 0,
    Reader,
    Writer,
    Editor
}

export const DocRightOptions: ListOptionsMap = {
    0: {
        val:DocRights.NoRight,
        str:"No right"
    },
    1: {
        val:DocRights.Reader,
        str:"Reader"
    },
    2: {
        val:DocRights.Writer,
        str:"Writer"
    },
    3: {
        val:DocRights.Editor,
        str:"Editor"
    },
}

export class Access extends Searchable{
    ID: number
    public_key: string
    title: string
    subtitle?: string
    isSelected: boolean
    right: Number
    isUserAccess: Boolean
    docID: Number

    static readonly Rights = DocRights;
    readonly Rights = Access.Rights;

    constructor(data: any, isUserAccess = true, docID = -1){
        super()
        this.ID = data.ID
        this.public_key = data.public_key
        this.title = data.name
        this.isSelected = false 
        this.isUserAccess = isUserAccess
        this.docID = docID
        if(data.email !== undefined){
            this.subtitle = data.email
        } else {
            this.subtitle = ""
        }
        this.right = data.right
        if(this.right == undefined){
            this.right = 0
        }
    }

    get repr(){
        var repr = this.title
        if(this.subtitle){
            repr += "- " + this.subtitle
        }
        return repr
    }

    get uri(){
        return "/api/docs/"+ this.docID + (this.isUserAccess ? "/users/" : "/groups/") + this.ID
    }

    payload(key: string):any {
        return {
            ID:this.ID,
            key: Crypto.encryptRSA(key,this.public_key),
            right: this.right
        }
    }

    search(text: string) {
        return (this.title+this.subtitle).toLowerCase().includes(text)
    }
}

export class Content{
    ID?: number
    content: string
    _content?: string
    iv: string

    constructor(content?: any){
        if(content !== undefined){
            this.ID = content.ID
            this.content = content.content
            this.iv = content.iv
        } else {
            this.iv = Crypto.genIV();
        }
    }

    get payload(): any {
        return {
            ID: this.ID,
            content: this.content,
            iv: this.iv
        }
    }

    getData(key: string): string {
        if (this._content === undefined){
            this._content = Crypto.decryptAES(this.content, key, this.iv).data
        }
        return this._content
    }
    
    setData(data: string, key:string) {
        this._content = data
        this.content = Crypto.encryptAES(data, key, this.iv).data
    }
    
}

export class Part extends Searchable {
    ID: number
    level: number
    order: number
    _title: Content
    _characteristics: Content
    _content: Content
    CreatedAt: string
    DeletedAt: string
    UpdatedAt: string
    doc: Doc
    isDeleted: boolean
    attachments: Record<string,Blob>
    attachmentUri = "attachments"
    
    static fromData(data: any[], doc: Doc): Part[]{
        var parts: Part[];
        parts = []
        data.forEach((e) => {parts.push(new Part(doc, e))})
        return parts
    }

    constructor(doc: Doc, data?: any){
        super()
        this.doc = doc
        this.attachments = {}
        if(data !== undefined){
            this.ID = data.ID
            this.CreatedAt = data.CreatedAt
            this.UpdatedAt = data.UpdatedAt
            this.DeletedAt = data.DeletedAt
            this.level = data.level
            this.order = data.order
            this._title = new Content(data.title)
            this._characteristics = new Content(data.characteristics)
            this._content = new Content(data.content)
        } else {
            this._title = new Content()
            this._characteristics = new Content()
            this._content = new Content()
        }
    }

    get title(): string {
        return this._title.getData(this.doc.key)
    }

    set title(title: string) {
        this._title.setData(title, this.doc.key)
    }

    get characteristics(): string {
        return this._characteristics.getData(this.doc.key)
    }

    set characteristics(title: string) {
        this._characteristics.setData(title, this.doc.key)
    }

    get content(): string {
        return this._content.getData(this.doc.key)
    }

    set content(title: string) {
        this._content.setData(title, this.doc.key)
    }

    get payload(): any{
        return {
            ID:this.ID,
            level: this.level,
            order: this.order,
            title: this._title.payload,
            characteristics: this._characteristics.payload,
            content: this._content.payload,
        }
    }
    
    get titlePayload(): any{
        return {
            ID:this.ID,
            level: this.level,
            order: this.order,
            title: this._title.payload,
        }
    }

    get contentPayload(): any{
        return {
            ID:this.ID,
            level: this.level,
            order: this.order,
            content: this._content.payload,
        }
    }

    get characteristicsPayload(): any{
        return {
            ID:this.ID,
            level: this.level,
            order: this.order,
            characteristics: this._characteristics.payload,
        }
    }

    get lightPayload(): any{
        return {
            ID:this.ID,
            level: this.level,
            order: this.order,
        }
    }

    get localAttachments(){
        const localAttachments = [
            ...findAttachmentIDs(this.content,this.attachmentUri),
            ...findAttachmentIDs(this.characteristics,this.attachmentUri)
        ]
        return localAttachments
    }

    get uri(){
        return `/api/docs/${this.doc.ID}/parts/${this.ID}`
    }

    get attachmentsApiUri(){
        return `/api/docs/${this.doc.ID}/parts/${this.ID}/attachments`
    }
    
    get wsUrl(){
        return `${getWSServerUrl()}/ws/docs/${this.doc.ID}/parts/${this.ID}`
    }

    get maxOrder(){
        return this.doc.parts.length - 1
    }

    verifyAttachments(){
        const localAttachments = this.localAttachments
        var attachmentsToDelete = []
        const {get, del} = useApiToken(this.attachmentsApiUri)

        get().then((data) => {
            data.forEach(element => {
                if(!(localAttachments.includes(element.ID))){
                    attachmentsToDelete.push(element.ID)
                }
            });
            if(attachmentsToDelete.length != 0){
                del({"data":attachmentsToDelete})
            }
        })
    }

    search(text:string){
        return (this.title).toLowerCase().includes(text)
    }

    orderUp(){
        if (this.order < this.doc.parts.length){
            //index == order
            let nextPart = this.doc.parts[this.order + 1]
            nextPart.order -=1
            nextPart.silentOrderSave()
            this.order +=1
            this.silentOrderSave()
            this.doc._sortParts()
        }
    }

    orderDown(){
        if (this.order > 0){
            //index == order
            let previousPart = this.doc.parts[this.order - 1]
            previousPart.order +=1
            previousPart.silentOrderSave()
            this.order -=1
            this.silentOrderSave()
            this.doc._sortParts()
        }
    }

    levelUp(){
        if (this.level < 3){
            this.level += 1
            this.silentOrderSave()
        }
    }

    levelDown(){
        if (this.level > 0){
            this.level -= 1
            this.silentOrderSave()
        }
    }

    export(attachmentsZip, next){
        var self = this
        var countAttachments = this.localAttachments.length
        var md = ""
        this.localAttachments.forEach(e => {
            self.getAttachment(e.toString(), (blob)=>{
                attachmentsZip.file(e,blob)
                countAttachments -= 1
            })
        })
        Object.entries(this.attachments).forEach(([k,v]) => attachmentsZip.file(k,v));
        function wait(){
            if(countAttachments != 0){
                console.log("wait for attachments",countAttachments)
                setTimeout(wait,50)
            } else {
                next(md)
            }
        }
        md = "#".repeat(this.level+1) + " " + this.title + "\n" 
        if (this.characteristics != ""){
            md += "#".repeat(this.level+2) + " Properties" + this.characteristics + "\n"

        }
        md += this.content.replaceAll(/^(?=#+)/gm,"#".repeat(this.level+1)) + "\n-----\n"
        wait()
    }
    
    save(){
        const {
            loading,
            post,
            error
        } = useApiToken(this.uri)
        const self = this
        const save = (payload=this.payload) =>{
            self.verifyAttachments()
            return new Promise<any>((resolve, reject) => {
                post(payload).then(()=> {
                    if (error.value) reject(error)
                    else {
                        console.log("Part is saved")
                        resolve(true)
                    } 
                })
            })
        }

        return {
            loading,
            error,
            save
        }
    }

    silentSave(payload=this.payload){
        this.save().save(payload)
    }

    silentOrderSave(){
        this.save().save(this.lightPayload)
    }

    delete(){
        const {
            loading,
            del: _del,
            error
        } = useApiToken(this.uri)

        const del = () =>{
            return new Promise<any>((resolve, reject) => {
                _del().then(()=> {
                    if (error.value) reject(error)
                    else {
                        this.isDeleted = true
                        this.doc.verifyOrder(this.order)
                        this.attachments = {}
                        resolve(true)
                    } 
                })
            })
        }

        return {
            loading,
            error,
            del
        }
    }

    silentDelete(){
        this.delete().del()
    }

    addAttachment(attachment: Blob|File){
        const {
            loading,
            post,
            error
        } = useApiToken(this.uri+"/attachments")

        const addAttachment = () =>{
            return new Promise<any>((resolve, reject) => {
                Crypto.encryptBlobAES(attachment, this.doc.key).then((encBlob: Blob)=>{
                    const attachmentData = new FormData()
                    attachmentData.append('attachment',encBlob)
    
                    post(attachmentData, {'headers':{'Content-type':'multipart/form-data'}}).then((data)=>{
                        if (error.value) reject(error)
                        else {
                            resolve(`${this.attachmentUri}/${data}`)
                        } 
                    })
                })
            })
        }
        
        return {
            loading,
            error,
            addAttachment
        }
    }

    getAttachment(uri, callback){
        const attachmentID = uri.split('/').slice(-1)[0]
        if(this.attachments[attachmentID] != undefined){
            callback(this.attachments[attachmentID])
        } else {
            const {get} = useApiToken(this.attachmentsApiUri+"/"+attachmentID)
    
            get().then((data)=>{
                this.attachments[attachmentID] = Crypto.decryptBlobAES(data, this.doc.key)
                callback(this.attachments[attachmentID])
            })
        }
    }
}

export class Doc extends Searchable {
    static readonly Rights = DocRights;
    readonly Rights = Access.Rights;

    ID: number
    _title: Content
    key: string
    _key: string
    right: number
    parts: Part[] = []
    CreatedAt: string
    DeletedAt: string
    UpdatedAt: string

    constructor(data: any){
        super()
        this.ID = data.ID
        this.CreatedAt = data.CreatedAt
        this.UpdatedAt = data.UpdatedAt
        this.DeletedAt = data.DeletedAt
        this._title = new Content(data.title)
        this._key = data.key
        this.right = data.right
        this.key = Crypto.decryptRSA(this._key)
        if(data.parts != null){
            data.parts.forEach(e => this.parts.push(new Part(this, e)));
            this._sortParts()
        }
    }

    get title(): string {
        return this._title.getData(this.key)
    }

    set title(title: string) {
        this._title.setData(title, this.key)
    }

    get payload(): any{
        return {
            ID:this.ID,
            title: this._title.payload,
        }
    }

    get uri(){
        return `/api/docs/${this.ID}`
    }

    _sortParts(){
        this.parts.sort((a,b)=>{
            if(a.order>b.order){
                return 1
            }
            if(a.order<b.order){
                return -1
            }
            return 0
        })
    }

    addPart(part: Part) {
        part.level = 0
        part.order = this.parts.length
        this.parts.push(part)
        this._sortParts()
    }

    movePart(start: number, end: number){
        // case 1: the part is moving from order i to i + j where j > 0
        let step = 1
        if (end < start){
            // case 2: the part is moving from order i to i - j where j > 0
            step = -1
        }
        this.parts[start].order = this.parts[end].order
        this.parts[start].silentOrderSave()

        for (let i = start+step; i != end+step; i+=step) {
            this.parts[i].order -= step
            this.parts[i].silentOrderSave()
        }
        this._sortParts()
        if(import.meta.env.DEV)
            this.parts.forEach(e => console.log(e.title, e.order))
    }

    verifyOrder(start=0, end = this.parts.length){
        var newOrder = start
        var indexToDelete = []
        for (let i = start; i < end; i++) {
            const part = this.parts[i];
            if(!part.isDeleted){
                if(newOrder != part.order){
                    part.order = newOrder
                    part.silentOrderSave()// Todo improve this to have only one request
                }
                newOrder += 1
            } else {
                indexToDelete.push(i)
            }
        }
        indexToDelete.forEach(i => this.parts.splice(i,1));
    }

    checkRight(rightToCheck: DocRights){
        return this.right >= rightToCheck ? true : false; 
    }

    encryptUpdate(data: Uint8Array): string{
        const encryptedData = Crypto.encryptAES(data, this.key)
        return btoa(encryptedData.iv + encryptedData.data)
    }

    decryptUpdate(data: string): Uint8Array{
        data = atob(data)
        const decryptedData = Crypto.decryptAES( data.substring(24), this.key,data.substring(0,24)).data
        let stateOrUpdate = new Uint8Array(decryptedData.length);
        for (let index = 0; index < decryptedData.length; index++) {
            stateOrUpdate[index] = decryptedData.charCodeAt(index);
            
        }
        return stateOrUpdate
    }

    search(text: string){
        return (this.title).toLowerCase().includes(text)
    }

    export(){
        var zip = new JSZip();
        var attachments = zip.folder("attachments")
        var docStr = `# Properties\n- title: ${this.title}\n\n`;
        var countParts = this.parts.length
        var partsStr = new Array(countParts)

        function partCallback(part){
            return (partStr) =>{
                partsStr[part.order] = partStr;
                countParts -= 1

                if(countParts == 0){
                    partsStr.forEach(e => docStr += e)
                    zip.file("document.md",docStr)
                    zip.generateAsync({type:"blob"})
                    .then(function(content) {
                        saveAs(content, "document.zip");
                    });
                }
            }
        }
        this.parts.forEach(e => e.export(attachments,partCallback(e)))
    }

    refreshParts(){
        const {
            loading,
            get,
            error
        } = useApiToken("/api/docs/"+this.ID+"/parts")

        const refreshParts = () =>{
            return new Promise<any>((resolve, reject) => {
                get().then((data)=>{
                    if (error.value) reject(error)
                    else {
                        this.parts = Part.fromData(data,this)
                        this._sortParts()
                        resolve(true)
                    } 
                })
            })
        }

        return {
            loading,
            error,
            refreshParts
        }
    }

    addAttachment(attachment: Blob|File){
        const {
            loading,
            data,
            post,
            error
        } = useApiToken("/api/docs/"+this.ID+"/attachments")

        const addAttachment = () =>{
            return new Promise<any>((resolve, reject) => {
                Crypto.encryptBlobAES(attachment, this.key).then((encBlob: Blob)=>{
                    const attachmentData = new FormData()
                    attachmentData.append('attachment',encBlob)
    
                    post(attachmentData, {'headers':{'Content-type':'multipart/form-data'}}).then((data)=>{
                        if (error.value) reject(error)
                        else {

                            resolve(`/api/docs/${this.ID}/attachments/${data}`)
                        } 
                    })
                })
            })
        }
        
        return {
            loading,
            error,
            addAttachment
        }
    }
}
