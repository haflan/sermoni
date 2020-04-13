<template>
    <div class="services-wrapper">
        <div v-for="service in services">
            {{ service.id }} : 
            <input type="text" :value="service.token">
            <input type="text" :value="service.name">
            <input type="text" :value="service.description">
            <input type="text" :value="service.period">
            <input type="text" :value="service.maxevents">
        </div>

        <input type="text" v-model="newService.token" placeholder="Token">
        <input type="text" v-model="newService.name" placeholder="Name">
        <input type="text" v-model="newService.description" 
               placeholder="Description">
        <input type="text" v-model="newService.period"
               placeholder="Expectation Period">
        <input type="text" v-model="newService.maxevents" 
               placeholder="Max number of events">

        <button @click="addService">Add service</button>
    </div>
</template>

<script>
    import api from "./requests.js";
    export default {
        name: "Services",
        data() {
            return {
                services: [],
                newService: {
                    token: "",
                    name: "",
                    description: "",
                    period: 0,
                    maxevents: 0
                }
            }
        },
        methods: {
            addService() {
                api.postService(this.newService,
                    success => {
                        this.services.push(this.newService);
                        this.newService = {
                            token: "", name: "", description: "",
                            period: 0, maxevents: 0
                        };
                    },
                    error => {
                        this.$emit("error");
                    }
                );
            }
        },
        mounted() {
            api.getServices(
                services => {
                    this.services = services;
                },
                error => {
                    console.log(error);
                    this.$emit("error");
                }
            );
        }
    }
</script>

<style scoped>
</style>
