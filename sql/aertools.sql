--
-- Database: `aertools`
--

-- --------------------------------------------------------

--
-- Table structure for table `accounts`
--

CREATE TABLE IF NOT EXISTS `accounts` (
`id` int(11) NOT NULL,
  `login` varchar(16) COLLATE utf8_unicode_ci NOT NULL,
  `level` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `auth_tokens`
--

CREATE TABLE IF NOT EXISTS `auth_tokens` (
  `user_id` int(11) NOT NULL,
  `token` varchar(64) COLLATE utf8_unicode_ci NOT NULL,
  `expiration` date NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `lockers`
--

CREATE TABLE IF NOT EXISTS `lockers` (
`id` int(11) NOT NULL,
  `login` varchar(16) COLLATE utf8_unicode_ci NOT NULL,
  `locker` varchar(8) COLLATE utf8_unicode_ci NOT NULL,
  `borrowing` int(11) NOT NULL,
  `retrieval` int(11) DEFAULT NULL,
  `state` varchar(64) COLLATE utf8_unicode_ci DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

--
-- Indexes for table `accounts`
--
ALTER TABLE `accounts`
 ADD PRIMARY KEY (`id`), ADD UNIQUE KEY `login` (`login`);

--
-- Indexes for table `auth_tokens`
--
ALTER TABLE `auth_tokens`
 ADD PRIMARY KEY (`user_id`), ADD UNIQUE KEY `token` (`token`);

--
-- Indexes for table `lockers`
--
ALTER TABLE `lockers`
 ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for table `accounts`
--
ALTER TABLE `accounts`
MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `lockers`
--
ALTER TABLE `lockers`
MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;