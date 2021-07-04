<template>
    <div class="flex flex-row h-full">
        <div class="flex-1 h-full flex flex-col overflow-hidden" v-for="accessList in accessLists" :class="{'mr-2': accessList.list == users}">
            <div class="my-2">
                <h1 class="text-grey-darkest mb-2">User Access</h1>
                <SearchInput :data="accessList.list" @found="(found) => accessList.filtered = found" />
            </div>
            <div class="flex-auto h-40 flex flex-col rounded rounded-b-none bg-secondary p-2 overflow-y-auto -pt-1 -pb-1">
                <div class="flex flex-row bg-primary p-3 rounded items-center my-1" v-for="element in accessList.filtered">
                    <span class="flex-1 no-selection">{{element.repr}}</span>
                    <List 
                        :modelValue="element.right"
                        :button="DocRightOptions[element.right.valueOf()].str"
                        :options="DocRightOptions"
                        @update:modelValue="updateRight(element, $event, putUserDoc)"
                    />
                    <input class="h-6 w-12 rounded" type=checkbox @click.stop="updateRight(element, element.right == 0 ? 1 : 0, accessList.put)" v-model="element.isSelected"/>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive } from "vue";
import { Access, Doc, DocRightOptions } from "@modules/models";
import { useApiToken } from '@modules/api';
import List from "@components/ui/List.vue";
import SearchInput from "@components/ui/SearchBar.vue";

interface IAccessListView {
    list: Access[];
    filtered: Access[];
    put: Function;
}

const props = defineProps({
    doc: Doc
})
const users = ref<Access[]>([])
const filteredUsers = ref<Access[]>([])
const groups = ref<Access[]>([])
const filteredGroups = ref<Access[]>([])
const accessLists = reactive<IAccessListView[]>([])

const {
    loading: loadUsersDoc,
    data: dataUsersDoc,
    get: getUsersDoc, put: putUserDoc,
    error: errorUsersDoc
} = useApiToken("/api/docs/"+props.doc.ID+"/users")

const {
    loading: loadGroupsDoc,
    data: dataGroupsDoc,
    get: getGroupsDoc, put: putGroupDoc,
    error: errorGroupsDoc
} = useApiToken("/api/docs/"+props.doc.ID+"/groups")

const {
    loading: loadUsers,
    data: dataUsers,
    get: getUsers,
    error: errorUsers
} = useApiToken("/api/users/verified")

const {
    loading: loadGroups,
    data: dataGroups,
    get: getGroups,
    error: errorGroups
} = useApiToken("/api/groups")

const loading = computed(() => loadUsersDoc || loadGroupsDoc || loadUsers || loadGroups)
const error = computed(() => errorUsersDoc || errorGroupsDoc || errorUsers || errorGroups)

accessLists.push({list:users.value, filtered: filteredUsers.value, put: putUserDoc})
accessLists.push({list:groups.value, filtered: filteredGroups.value, put: putGroupDoc})

function initAccessesData(accesses, data, isUser = true): Function{
    return (e) => {
        var access = new Access(e, isUser, props.doc.ID)
        let dataFound = data.value.find((e) => e.ID == access.ID)
        if(dataFound != undefined){
            access.isSelected = true
            access.right = dataFound.right
        }
        accesses.value.push(access)
    }
}

function sortData(e1,e2){
    if (e1.right > e2.right)
        return -1;
    if (e1.right < e2.right)
        return 1;
    return 0;
}

getUsers().then(() => {
    getUsersDoc().then(()=>{
        dataUsers.value.forEach(initAccessesData(users, dataUsersDoc))
        users.value.sort(sortData)
    });
});

getGroups().then(() => {
    getGroupsDoc().then(()=>{
        dataGroups.value.forEach(initAccessesData(groups, dataGroupsDoc, false))
        groups.value.sort(sortData)
    });
});

function updateRight(access: Access, newRight, put: Function){
    access.right = newRight;
    if(access.right != DocRightOptions[0].val){
        if(access.isSelected){
            const { post } = useApiToken(access.uri)
            post({right: newRight})
        } else {
            put(access.payload(props.doc.key)).then((data)=>{
                if(data == 0){
                    access.isSelected = true
                }
            })
        }
    } else {
        const { del } = useApiToken(access.uri)
        del().then(()=>{access.isSelected = false})
    }
}
</script>

<style lang="scss">

</style>