<template>
  <div class="bg-gray-800 p-6 rounded-lg shadow-md mb-8">
    <h2 class="text-2xl font-bold mb-4">最活跃日小时分布 🕒</h2>
    <div class="h-64">
      <Line :data="chartData" :options="chartOptions" />
    </div>
  </div>
</template>

<script>
import { computed } from 'vue';
import { Line } from 'vue-chartjs';
import { Chart as ChartJS, Title, Tooltip, Legend, LineElement, PointElement, CategoryScale, LinearScale } from 'chart.js';

ChartJS.register(Title, Tooltip, Legend, LineElement, PointElement, CategoryScale, LinearScale);

export default {
  name: 'HourlyVisitChart',
  components: { Line },
  props: {
    hourlyDistribution: Object
  },
  setup(props) {
    const chartData = computed(() => {
      const sortedHours = Object.entries(props.hourlyDistribution || {})
          .sort(([a], [b]) => parseInt(a) - parseInt(b));
      return {
        labels: sortedHours.map(([hour]) => `${hour}:00`),
        datasets: [{
          label: '小时访问分布',
          borderColor: 'rgba(255, 99, 132, 1)',
          data: sortedHours.map(([, count]) => count),
          tension: 0.1
        }]
      };
    });

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
