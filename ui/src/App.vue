<template>
    <div>
        <header> 
            <div id="bar">
                <div style="font-size: 1.5em; color: #bbf">&gt; sermoni</div> 
                <!-- color #bbf instead? -->
                <div style="margin-left: auto;">
                    <img src="dist/eye.png" 
                         @click="togglePage"
                         style="height: 3em; opacity: 0.3;">
                </div>
            </div>
        </header>
        <main>
            <component :is="page"/>
        </main>
    </div>
</template>

<script>
    import Events from "./Events.vue";
    import Services from "./Services.vue";
    export default {
        name: "App",
        components: {Events, Services},
        data() {
            return {
                page: Events,
            };
        },
        methods: {
            togglePage() {
                // Should do nothing when on login page
                if (this.page === Events) {
                    this.page = Services;
                } else if (this.page === Services) {
                    this.page = Events;
                }
            }
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
    background-color: #eef;
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
