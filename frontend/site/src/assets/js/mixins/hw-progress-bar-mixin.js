/* 
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */

let HwProgressBarMixin = {

    data: function () {
        return {}
    },

    created: function () {

    },

    methods: {
        fillRamProgressBar(item_id, min, value, max) {
            let percent_colors = [
                {percent: min / max, color: {r: 0xff, g: 0x00, b: 0}},
                {percent: value / max, color: {r: 0xff, g: 0xff, b: 0xff}},
                {percent: 1.0, color: {r: 0xff, g: 0xff, b: 0xff}}];

            this.fillProgressBar(`item-${item_id}-ram-progress`, percent_colors);
        },
        fillFreeRamProgressBar(item_id, min, value, max) {
            let percent_colors = [
                {percent: min / max, color: {r: 0xff, g: 0xa4, b: 0}},
                {percent: value / max, color: {r: 0xff, g: 0xff, b: 0xff}},
                {percent: 1.0, color: {r: 0xff, g: 0xff, b: 0xff}}];

            this.fillProgressBar(`item-${item_id}-free-ram-progress`, percent_colors);
        },
        fillDiskProgressBar(item_id, min, value, max) {
            let percent_colors = [
                {percent: min / max, color: {r: 0xff, g: 0xff, b: 0}},
                {percent: value / max, color: {r: 0xff, g: 0xff, b: 0xff}},
                {percent: 1.0, color: {r: 0xff, g: 0xff, b: 0xff}}];

            this.fillProgressBar(`item-${item_id}-disk-progress`, percent_colors);
        },
        fillCpuProgressBar(item_id, min, value, max) {
            let percent_colors = [
                {percent: min / max, color: {r: 0x00, g: 0xff, b: 0x00}},
                {percent: value / max, color: {r: 0xff, g: 0xff, b: 0xff}},
                {percent: 1.0, color: {r: 0xff, g: 0xff, b: 0xff}}];

            this.fillProgressBar(`item-${item_id}-cpu-progress`, percent_colors);
        },
        fillProgressBar(progress_bar_id, percent_colors) {
            const percents_total = 25;
            let progress_bar = document.getElementById(progress_bar_id);
            progress_bar.innerText = '';

            for (let i = 0, length = percents_total; i <= length; i++) {
                let li = document.createElement('li');
                let percent = i / length;
                li.style.backgroundColor = this.getColorForPercentage(percent_colors, percent);
                progress_bar.appendChild(li);
            }
        },
        getColorForPercentage(percent_colors, percent) {
            let i;

            for (i = 1; i < percent_colors.length - 1; i++) {
                if (percent < percent_colors[i].percent) {
                    break;
                }
            }
            let lower = percent_colors[i - 1];
            let upper = percent_colors[i];
            let range = upper.percent - lower.percent;
            let range_percent = (percent - lower.percent) / range;
            let percent_lower = 1 - range_percent;
            let percent_upper = range_percent;
            let color = {
                r: Math.floor(lower.color.r * percent_lower + upper.color.r * percent_upper),
                g: Math.floor(lower.color.g * percent_lower + upper.color.g * percent_upper),
                b: Math.floor(lower.color.b * percent_lower + upper.color.b * percent_upper)
            };
            return 'rgb(' + [color.r, color.g, color.b].join(',') + ')';
        }

    }

}

export {HwProgressBarMixin};