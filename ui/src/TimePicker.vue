<template>
    <div style="display: flex;">
        <input type="number" @input="update('number', $event)" :value="number"/>
        <select @input="update('scalar', $event)" :value="scalar">
            <option v-for="u in units" 
                    :key="u.unit"
                    :value="u.scalar" 
                    :selected="scalar === u.scalar">
                {{u.unit}}
            </option>
        </select>
    </div>
</template>

<script>
    export default {
        name: "TimePicker",
        props: {
            value: Object
        },
        data() {
            return {
                number: this.value.number,
                scalar: this.value.scalar,
                units:  [
                    { "unit": "weeks", "scalar": 604800000 },
                    { "unit": "days", "scalar": 86400000 },
                    { "unit": "hours", "scalar": 3600000 },
                    { "unit": "minutes", "scalar": 60000 },
                    { "unit": "seconds", "scalar": 1000 },
                    { "unit": "none", "scalar": 0 },
                ]
            }
        },
        watch: {
            value: {
                immediate: true,
                handler(val) {
                    this.number = val.number;
                    this.scalar = val.scalar;
                }
            }
        },
        methods: {
            update(prop, e) {
                const newValue = this.value;
                newValue[prop] = e.target.value;
                this.$emit("input", newValue);
            },

        }
    }
</script>