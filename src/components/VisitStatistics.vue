<template>
  <div class="bg-gray-800 p-6 rounded-lg shadow-md mb-8">
    <h2 class="text-2xl font-bold mb-4">访问统计 📊</h2>
    <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-6">
      <StatItem title="总访问次数" :value="profile.total_visits" bgColor="bg-blue-900" />
      <StatItem title="正常访问" :value="profile.total_normal_visits" bgColor="bg-green-900" />
      <StatItem title="攻击次数" :value="profile.total_attacks" bgColor="bg-red-900" />
      <StatItem title="最活跃日期" :value="profile.most_active_day" bgColor="bg-yellow-900" />
      <StatItem title="最常访问URL" :value="profile.most_visited_url" bgColor="bg-purple-900" />
    </div>
  </div>
</template>

<script>
import StatItem from './StatItem.vue';

export default {
  name: 'VisitStatistics',
  components: {
    StatItem
  },
  props: {
    profile: Object
  }
};
</script>
