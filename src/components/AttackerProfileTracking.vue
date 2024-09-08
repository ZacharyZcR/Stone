<template>
  <div class="bg-gray-900 text-white flex flex-col min-h-screen">
    <HeaderPage />

    <div class="container mx-auto px-4 py-8 flex-1 mt-16">
      <AttackerInfo :ip-address="ipAddress" @fetch-profile="fetchAttackerProfile" />
      <VisitStatistics v-if="profile" :profile="profile" />
      <DailyVisitChart v-if="profile" :daily-attacks="profile.daily_attacks" />
      <HourlyVisitChart v-if="profile" :hourly-distribution="profile.hourly_distribution" />
    </div>

    <FooterPage />
  </div>
</template>

<script>
import { ref } from 'vue';
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
    const ipAddress = ref('');
    const profile = ref(null);

    const fetchAttackerProfile = async (ip) => {
      try {
        const response = await api.get(`/attacker-profile?ip=${ip}`);
        profile.value = response.data;
      } catch (error) {
        console.error('获取攻击者画像失败:', error);
      }
    };

    return {
      ipAddress,
      profile,
      fetchAttackerProfile
    };
  }
};
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Roboto:wght@400;700&display=swap');
</style>
