<template>
    <div class="services-wrapper">
        <div class="service" v-for="service in services" :key="service.id">
            <span>Service ID:</span>
            <input type="text" :value="service.id"> <br/>

            <span>Token:</span>
            <input :type="showPasswords ? 'text' : 'password'" :value="service.token"/> <br/>

            <span>Name:</span>
            <input type="text" :value="service.name"/> <br/>

            <span>Description:</span>
            <input type="text" :value="service.description"/> <br/>

            <span>Max number of events:</span>
            <input type="number" :value="service.maxevents"/> <br/>

            <span>Expectation period:</span>
            <time-picker :value="service.period"/> <br/>

            <button @click="deletionID = service.id">Delete</button>
            <button @click="updateService(service.id)">Update</button>
        </div>

        <input :type="showPasswords ? 'text' : 'password'" v-model="newService.token" placeholder="Token"> <br/>
        <input type="text" v-model="newService.name" placeholder="Name"> <br/>
        <input type="text" v-model="newService.description" placeholder="Description"> <br/>
        <input type="number" v-model.number="newService.maxevents" placeholder="Max number of events"> <br/>
        <time-picker v-model="newService.period" placeholder="Expectation Period"/> <br/>

        <button @click="addService">Add service</button>
        <div v-show="deletionID" style="position: fixed; bottom: 15px; right: 15px;">
            <button @click="deletionID = 0">Cancel</button>
            <button @click="deleteService()">Confirm</button>
        </div>
    </div>
</template>

<script>
    import api from "./requests.js";
    import TimePicker from "./TimePicker.vue";
    export default {
        name: "Services",
        components: {TimePicker},
        data() {
            return {
                services: [],
                newService: {
                    token: "",
                    name: "",
                    description: "",
                    period: {"number": 0, "scalar": 0},
                    maxevents: 0
                },
                showPasswords: false,
                deletionID: 0,
            }
        },
        methods: {
            addService() {
                // Make unix milli time formatted period before sending
                const formattedService = {
                    ...this.newService,
                };
                formattedService.period = this.newService.period.scalar * this.newService.period.number;
                console.log(formattedService);
                api.postService(formattedService,
                    success => {
                        this.services.push(this.newService);
                        this.newService = {
                            token: "", name: "", description: "",
                            period: {"number": 0, "scalar": 0}, maxevents: 0
                        };
                    },
                    error => {
                        this.$emit("error");
                    }
                );
            },
            updateService(id) {
                alert("Not implemented!");
            },
            deleteService() {
                api.deleteService(this.deletionID,
                    success => {
                        this.services = this.services.filter(
                            s => s.id !== this.deletionID
                        );
                        this.deletionID = 0;
                    },
                    error => {
                        console.error(error);
                        this.$emit("error");
                    }
                );
            },
            getExpectations(unixMilliTime) {
                const units = [
                    { "unit": "weeks", "scalar": 604800000 },
                    { "unit": "days", "scalar": 86400000 },
                    { "unit": "hours", "scalar": 3600000 },
                    { "unit": "minutes", "scalar": 60000 },
                    { "unit": "seconds", "scalar": 1000 },
                    { "unit": "none", "scalar": 0 }
                ];
                if (unixMilliTime === 0) {
                    return { "number": 0, "scalar": 0 };
                }
                for (let i = 0; i < units.length; i++) {
                    let u = units[i];
                    if (unixMilliTime % u.scalar === 0) {
                        return {
                            "scalar": u.scalar,
                            "number": unixMilliTime / u.scalar
                        };
                    }
                }
            }
        },
        mounted() {
            api.getServices(
                services => {
                    this.services = services.map(s => {
                        s.period = this.getExpectations(s.period);
                        return s;
                    });
                },
                error => {
                    console.log(error);
                    this.$emit("error");
                }
            );
        }
    }
</script>

<style >
span {
    font-weight: bold;
}
input {
    border-width: 0 0 1px;
    border-color: #bbf;
    background-color: inherit;
    outline-width: 0;
    flex: 1;
    margin-bottom: 3px;
    width: 100%;
}
input:focus {
    outline-width: 0;
}
input::placeholder {
    color: #bbf;
}
div .service {
    border: 1px solid blueviolet;
    border-radius: 3px;
    margin-bottom: 1em;
    padding: 1em;
}
</style>
