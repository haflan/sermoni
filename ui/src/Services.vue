<template>
    <div class="services-wrapper">
        <div v-for="service in services">
            {{ service.name }} : 
            <input type="text" service.description }} : 
            <input type="text" :value="service.token">
        </div>
    </div>
</template>

<script>
    import api from "./requests.js";
    export default {
        name: "Services",
        data() {
            return {
                services: []
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
