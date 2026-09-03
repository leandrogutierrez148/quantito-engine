-- ============================================================================
-- QUANTITO ENGINE - Initial Execution Setup Data
-- Database: quantito-engine
-- ============================================================================

USE `quantito-engine`;

-- Disable foreign key checks during batch insert
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE `signals`;
TRUNCATE TABLE `variants`;
TRUNCATE TABLE `strategies`;
TRUNCATE TABLE `experiment_assets`;
TRUNCATE TABLE `experiments`;
SET FOREIGN_KEY_CHECKS = 1;

-- ============================================================================
-- 1. STRATEGIES (Quantitative Strategies Catalog)
-- ============================================================================
INSERT INTO `strategies` (`id`, `name`, `class`, `params`, `status`, `created_at`, `updated_at`) VALUES
(
  1,
  '50-Day Moving Average',
  'fifty_day_avg',
  JSON_OBJECT(
    'fifty_day_avg', 100.0,
    'current_price', 105.5,
    'threshold', 0.02
  ),
  'ACTIVE',
  NOW(),
  NOW()
),
(
  2,
  'Price Momentum',
  'price_momentum',
  JSON_OBJECT(
    'current_price', 150.0,
    'pct_change', 0.065,
    'buy_threshold', 0.05,
    'sell_threshold', -0.05
  ),
  'ACTIVE',
  NOW(),
  NOW()
),
(
  3,
  'Conservative Momentum',
  'price_momentum',
  JSON_OBJECT(
    'current_price', 150.0,
    'pct_change', 0.02,
    'buy_threshold', 0.08,
    'sell_threshold', -0.08
  ),
  'ACTIVE',
  NOW(),
  NOW()
);

-- ============================================================================
-- 2. EXPERIMENTS (Quantitative Experiments)
-- ============================================================================
INSERT INTO `experiments` (`id`, `experiment`, `description`, `params`, `status`, `notes`, `created_at`, `updated_at`) VALUES
(
  1,
  'Tech_Contraposition_SMA_vs_Momentum',
  'Estudio comparativo entre SMA de 50 dias y Price Momentum en acciones tecnologicas',
  JSON_OBJECT(
    'benchmark', 'SPY',
    'execution_mode', 'PAPER',
    'universe', 'US_TECH_EQUITIES'
  ),
  'ACTIVE',
  'Experimento de referencia para validacion y ejecucion del motor cuantitativo',
  NOW(),
  NOW()
),
(
  2,
  'Crypto_Momentum_Baseline',
  'Estudio de momentum en criptomonedas (Pausado)',
  JSON_OBJECT(
    'benchmark', 'BTC',
    'execution_mode', 'BACKTEST'
  ),
  'PAUSED',
  'Pausado para calibracion de umbrales',
  NOW(),
  NOW()
);

-- ============================================================================
-- 3. EXPERIMENT_ASSETS (Asset Assignment to Experiments)
-- ============================================================================
-- Experiment 1: Assets 101 (AAPL), 102 (MSFT), 103 (NVDA)
-- Experiment 2: Asset 201 (BTC)
INSERT INTO `experiment_assets` (`id`, `experiment_id`, `asset_id`, `status`, `created_at`, `updated_at`) VALUES
(1, 1, 101, 'ACTIVE', NOW(), NOW()),
(2, 1, 102, 'ACTIVE', NOW(), NOW()),
(3, 1, 103, 'ACTIVE', NOW(), NOW()),
(4, 2, 201, 'ACTIVE', NOW(), NOW());

-- ============================================================================
-- 4. VARIANTS (Experiment Variants / Strategy Configurations in Contraposition)
-- ============================================================================
-- Experiment 1: Contraposition between SMA 50 (Variant A) and Momentum (Variant B)
INSERT INTO `variants` (`id`, `experiment_id`, `strategy_id`, `variant`, `notes`, `created_at`, `updated_at`) VALUES
(
  1,
  1,
  1,
  'Variant A (50-Day SMA)',
  'Estrategia tendencial basada en media movil de 50 dias',
  NOW(),
  NOW()
),
(
  2,
  1,
  2,
  'Variant B (Price Momentum)',
  'Estrategia de momentum basada en variacion porcentual reciente',
  NOW(),
  NOW()
),
(
  3,
  2,
  3,
  'Variant Crypto Conservative',
  'Estrategia de momentum conservadora',
  NOW(),
  NOW()
);
