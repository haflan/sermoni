<template>
    <div>
        <header :style="headerStyle">
            <div id="bar">
                <div style="font-size: 1.5em; color: #bbf">&gt; sermoni</div> 
                <div @click="togglePage" style="margin-left: auto;">
                    <Eye :service-view="this.serviceView" style="cursor: pointer;"/>
                </div>
            </div>
        </header>
        <main>
            <component v-if="page"
                       :is="page" 
                       @login="login" 
                       @error="error = true"/>
        </main>
    </div>
</template>

<script>
    import Eye from "./Eye.vue";
    import Login from "./Login.vue";
    import Events from "./Events.vue";
    import Services from "./Services.vue";
    import api from "./requests.js";
    export default {
        name: "App",
        components: {Login, Eye, Events, Services},
        data() {
            return {
                page: null,
                serviceView: false,
                error: false
            };
        },
        methods: {
            login() {
                this.page = Events;
                this.error = false;
            },
            togglePage() {
                // Should do nothing when on login page
                if (this.page === Events) {
                    this.page = Services;
                    this.serviceView = true;
                } else if (this.page === Services) {
                    this.page = Events;
                    this.serviceView = false;
                }
            }
        },
        computed: {
            headerStyle() {
                const bgColor = this.error ? "#fce1e1" : "#eef";
                return {
                    "background-color": bgColor
                };
            }
        },
        mounted() {
            api.init(
                successData => {
                    if (successData.authenticated) {
                        this.page = Events;
                    } else {
                        this.page = Login;
                    }
                },
                errorData => {
                    console.error(errorData);
                    this.error = true;
                }
            )
        }
    }
</script>

<style>
body {
    font-family: Roboto, sans-serif;
    background-color: #FAFAFA;
    margin: 0;
}
header {
    z-index: 1000;
    border-bottom: 1px solid rgba(0,0,0,.075);
    -webkit-box-shadow: 0 0 10px rgba(0,0,0,.1);
    box-shadow: 0 0 10px rgba(0,0,0,.1);
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    padding: 0;
    height: 4em;
    display: flex;
    align-items: center; 
}
#bar {
    width: 100%;
    display: flex; 
    margin: 1em; 
    align-items: center; 
}
main {
    margin: 4em 0 0 0;
    width: calc(100%-1em);
    min-height: 1em;
    padding: 1em;
}
</style>
