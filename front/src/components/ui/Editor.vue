<template>
    <div class="z-0 overflow-y-auto h-full" ref="toastuiEditor" id="editor"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch  } from 'vue';

import { escapeXml, themeIsLight, usePropsWrapper } from "@modules/helpers";

import Editor from '@toast-ui/editor';
import chart from '@toast-ui/editor-plugin-chart';
import codeSyntaxHighlight from '@toast-ui/editor-plugin-code-syntax-highlight';
import colorSyntax from '@toast-ui/editor-plugin-color-syntax';
import uml from '@toast-ui/editor-plugin-uml';
// import tableMergedCell from '@toast-ui/editor-plugin-table-merged-cell';

import { getCollabWebSocket } from '@modules/collaboration';

const props = defineProps({
    modelValue: String,
    viewer: {type: Boolean, default: false},
    cwsID: {type: Number,default: -1},
})
const emit = defineEmits<{
    (event: 'update:modelValue', value: string)
    (event: 'uploadAttachment', attachment, callback)
    (event: 'loadImage', url: string,node: Element): void
    (event: 'save'): void
}>()

const value = usePropsWrapper(props, emit, "modelValue")
const toastuiEditor = ref<HTMLDivElement>();
var editor;
var saveTimeout: ReturnType<typeof setTimeout>;

function loadImage(id,url){
    function loadImageWhenReady(){
        var nodeList = document.querySelectorAll(`[data-nodeid="${id}"]`)
        if(nodeList.length != 0){
            var node = nodeList[0]
            emit("loadImage", url, node)
        } else {
            // While it must be asynchronous, 
            // we need load the image blob only when the dom node has been inserted.
            setTimeout(loadImageWhenReady, 10)
        }
    }
    loadImageWhenReady();
}

onMounted(()=>{
    var eventLoadHandler = (editor)=>{
        //@ts-ignore
        if(import.meta.env.DEV)
            console.log("editor.onload default")
        editor.setMarkdown(value.value, false);
    }
    var eventChangeHandler = () => {
        //@ts-ignore
        if(import.meta.env.DEV)
            console.log("editor.onchange default")
        clearTimeout(saveTimeout)
        value.value = editor.mdEditor.getMarkdown();
        saveTimeout = setTimeout(() => emit(`save`), 2000);
    }
    var cws = null
    var dontSaveThe1stUpdate = true; // Little trick to prevent 2nd+ clients from saving the document on upload.
    var plugins = [chart, uml, codeSyntaxHighlight, colorSyntax, uml]

    if(props.cwsID != -1){
        cws = getCollabWebSocket(props.cwsID)
        eventLoadHandler = (editor)=>{
            //@ts-ignore
            if(import.meta.env.DEV)
                console.log("editor.onload",cws.isFirstClient)
            if (cws.isFirstClient){
                editor.setMarkdown(value.value, false);
            }
            cws.editorIsLoaded(value.value);
        }
        eventChangeHandler = () => {
            //@ts-ignore
            if(import.meta.env.DEV)
                console.log("editor.onchange",cws.isFirstClient,cws.isInitialized, cws.mustSendUpdate)
            if(cws.isInitialized && cws.mustSendUpdate){
                if(!cws.isFirstClient && dontSaveThe1stUpdate){
                    dontSaveThe1stUpdate = false
                } else {
                    clearTimeout(saveTimeout)
                    value.value = editor.mdEditor.getMarkdown();
                    saveTimeout = setTimeout(() => emit(`save`), 2000);
                }
            }
        }
        plugins = plugins.concat([...cws.plugins])
        console.log(plugins)
    }
    editor = Editor.factory({
        el: toastuiEditor.value,
        height:'',
        //@ts-ignore
        viewer: props.viewer,
        initialEditType: 'markdown',
        usageStatistics: false,
        theme : themeIsLight.value ? 'light' : 'dark',
        plugins: plugins,
        events: {
            "load":eventLoadHandler,
            "change": eventChangeHandler
        }, 
        customHTMLRenderer: {
            image(node, context) {
                context.skipChildren()
                //@ts-ignore
                loadImage(node.id, node.destination!) // Handle the end-to-end encryption here
                return {
                    type: 'openTag',
                    tagName: 'img',
                    selfClose: true,
                    attributes: {
                        alt: context.getChildrenText(node),
                        //@ts-ignore
                        ...(node.title && { title: escapeXml(node.title) }),
                    },
                }
            }
        }
    });
    // hack to remove the default drop image handling as we define our handler
    if(!props.viewer){
        editor.removeHook('addImageBlobHook');
        editor.addHook('addImageBlobHook',(attachment,callback)=>{
            emit('uploadAttachment',attachment, callback)
        });
    }

    watch(themeIsLight, (newThemeIsLight)=>{
        if(newThemeIsLight){
            (toastuiEditor.value.firstChild as HTMLDivElement).classList.remove('toastui-editor-dark');
            (toastuiEditor.value.firstChild as HTMLDivElement).classList.add('toastui-editor');
        } else {
            (toastuiEditor.value.firstChild as HTMLDivElement).classList.add('toastui-editor-dark');
            (toastuiEditor.value.firstChild as HTMLDivElement).classList.remove('toastui-editor');
        }
    })
})


</script>

<style>
.toastui-editor-defaultUI{
    border-radius: 0px;
}

.toastui-editor-defaultUI .ProseMirror{
    height: 100%;
    margin-bottom: 10px;
}
    
.toastui-editor-defaultUI-toolbar{
    border-top-right-radius: 0px !important;
}

.toastui-editor-md-tab-container{
    border-top-left-radius: 0px !important;
}

.toastui-editor-md-tab-container .tab-item{
    margin-top: 12px;
}

.toastui-editor-mode-switch{
    display: none !important;
}
</style>