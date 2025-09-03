<template>
  <div class="bg-gray-900 text-white flex flex-col min-h-screen">
    <HeaderPage />

    <div class="container mx-auto px-4 py-8 flex-1 mt-16">
      <AttackerInfo :ip-address="ipAddress" @fetch-profile="fetchAttackerProfile" />
      <div v-if="loading" class="text-center py-8">
        <p class="text-xl">加载中...</p>
      </div>
      <div v-else-if="error" class="text-center py-8">
        <p class="text-xl text-red-500">{{ error }}</p>
      </div>
      <template v-else-if="profile">
        <VisitStatistics :profile="profile" />
        <DailyVisitChart :daily-attacks="profile.daily_attacks" />
        <HourlyVisitChart :hourly-distribution="profile.hourly_distribution" />
      </template>
      <div v-else class="text-center py-8">
        <p class="text-xl">没有找到数据</p>
      </div>
    </div>

    <FooterPage />
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import HeaderPage from './HeaderPage.vue';
import FooterPage from './FooterPage.vue';
import AttackerInfo from './AttackerInfo.vue';
import VisitStatistics from './VisitStatistics.vue';
import DailyVisitChart from './DailyVisitChart.vue';
import HourlyVisitChart from './HourlyVisitChart.vue';
import api from '../api/axiosInstance';

export default {
  name: 'AttackerProfileTracking',
  components: {
    HeaderPage,
    FooterPage,
    AttackerInfo,
    VisitStatistics,
    DailyVisitChart,
    HourlyVisitChart
  },
  setup() {
    const ipAddress = ref('127.0.0.1');
    const profile = ref(null);
    const loading = ref(false);
    const error = ref(null);

    const fetchAttackerProfile = async (ip) => {
      loading.value = true;
      error.value = null;
      try {
        const response = await api.get(`/attacker-profile?ip=${ip}`);
        profile.value = response.data;
        ipAddress.value = ip;
      } catch (err) {
        console.error('获取攻击者画像失败:', err);
        error.value = '获取数据失败，请稍后再试';
        profile.value = null;
      } finally {
        loading.value = false;
      }
    };

    onMounted(() => {
      fetchAttackerProfile(ipAddress.value);
    });

    return {
      ipAddress,
      profile,
      fetchAttackerProfile,
      loading,
      error
    };
  }
};
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Roboto:wght@400;700&display=swap');
</style>
