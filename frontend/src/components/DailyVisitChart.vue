<template>
  <div class="bg-gray-800 p-6 rounded-lg shadow-md mb-8">
    <h2 class="text-2xl font-bold mb-4">每日访问分布 📈</h2>
    <div class="h-64">
      <Bar :data="chartData" :options="chartOptions" />
    </div>
  </div>
</template>

<script>
import { computed } from 'vue';
import { Bar } from 'vue-chartjs';
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale } from 'chart.js';

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale);

export default {
  name: 'DailyVisitChart',
  components: { Bar },
  props: {
    dailyAttacks: Object
  },
  setup(props) {
    const chartData = computed(() => ({
      labels: Object.keys(props.dailyAttacks || {}),
      datasets: [
        {
          label: '总访问',
          backgroundColor: 'rgba(75, 192, 192, 0.6)',
          data: Object.values(props.dailyAttacks || {}).map(day => day.total)
        },
        {
          label: '正常访问',
          backgroundColor: 'rgba(54, 162, 235, 0.6)',
          data: Object.values(props.dailyAttacks || {}).map(day => day.normal)
        },
        {
          label: '攻击',
          backgroundColor: 'rgba(255, 99, 132, 0.6)',
          data: Object.values(props.dailyAttacks || {}).map(day => day.attacks)
        }
      ]
    }));

    const chartOptions = {
      responsive: true,
      maintainAspectRatio: false,
      scales: {
        y: {
          beginAtZero: true,
          ticks: { color: 'white' },
          grid: { color: 'rgba(255, 255, 255, 0.1)' }
        },
        x: {
          ticks: { color: 'white' },
          grid: { color: 'rgba(255, 255, 255, 0.1)' }
        }
      },
      plugins: {
        legend: {
          display: true,
          labels: { color: 'white' }
        },
        tooltip: {
          mode: 'index',
          intersect: false,
        }
      }
    };

    return { chartData, chartOptions };
  }
};
</script>
